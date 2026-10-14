import axios from "axios"
// yükleme animasyonu için
import { triggerLoader } from "../utils/eventBus";
import { toast } from "react-toastify";

const baseUrl = "https://ecommerceapi.yusufsevinc.work/api/";

export const api = axios.create({
    baseURL: baseUrl,
    withCredentials: true,
});

// Kuyruk ve Yenileme Durumunu Tutan Değişkenler
let isRefreshing = false;
let failedQueue: Array<{ resolve: (value?: unknown) => void, reject: (reason?: any) => void }> = [];

// Bekleyen İstekleri İşleyen (Kuyruğu Boşaltan) Fonksiyon
const processQueue = (error: any) => {
    failedQueue.forEach(promise => {
        if (error) {
            promise.reject(error);
        } else {
            promise.resolve();
        }
    });
    // İşlem bitince kuyruğu temizle
    failedQueue = [];
};

api.interceptors.request.use(
    (config) => {
        triggerLoader(true);
        return config;
    },
    (error) => {
        triggerLoader(false);
        return Promise.reject(error);
    }
);

api.interceptors.response.use(
    (response) => {
        console.log("interceptor çalıştı")
        triggerLoader(false);
        return response;
    },
    async (error) => {
        console.log("interceptor çalıştı error")
        triggerLoader(false);

        const originalRequest = error.config;

        if (error.response?.status === 401 && !originalRequest._retry) {
            originalRequest._retry = true;

            if (isRefreshing) {
                // Bu isteği yeni token alınana kadar Promise olarak kuyruğa (failedQueue) ekle
                return new Promise(function (resolve, reject) {
                    failedQueue.push({ resolve, reject });
                })
                .then(() => {
                    // Refresh işlemi başarılı olduysa bu orijinal isteği tekrar at
                    return api(originalRequest);
                })
                .catch((err) => {
                    // Refresh başarısız olduysa bu isteği de iptal et
                    return Promise.reject(err);
                });
            }

            isRefreshing = true;

            return new Promise(function (resolve, reject) {
                // Refresh token isteğini at
                axios.post(baseUrl + "refreshtoken", {}, {
                    withCredentials: true,
                    headers: { 'Content-Type': 'application/json' }
                })
                .then((response) => {
                    if (response.status === 200) {
                        // Token başarıyla yenilendi! Kuyrukta bekleyen tüm istekleri serbest bırak.
                        processQueue(null); 
                        // Kendi isteğimizi (orijinal isteği) de tekrar atarak sonucu döndür
                        resolve(api(originalRequest));
                    }
                })
                .catch((refreshError) => {
                    // Refresh token süresi de dolmuş veya geçersiz!
                    // Kuyrukta bekleyen diğer isteklere de "Hata var, iptal olun" mesajı gönder.
                    processQueue(refreshError); 
                    
                    if (localStorage.getItem("isAuthenticated") === "true") {
                        toast.error("Oturum süreniz doldu, lütfen tekrar giriş yapın.");
                    }
                    forceLogout();
                    
                    reject(refreshError);
                })
                .finally(() => {
                    // İşlem bitti (başarılı veya başarısız), bayrağı (flag) indir
                    isRefreshing = false;
                });
            });
        }

        // Catch 403 Forbidden from the backend
        if (error.response?.status === 403) {
            console.error("Yetkisiz Erişim! Sayfa zorla açılmaya çalışıldı, çıkış yapılıyor.");
            if (localStorage.getItem("isAuthenticated") === "true") {
                toast.error("Yetkisiz Erişim! Çıkış yapılıyor.");
            }
            forceLogout();
        } else if (error.response?.status !== 401) {
            // General error toast for all other errors except 401 which are handled by refresh
            const serverMessage = error.response?.data?.message || "Bir hata oluştu.";
            toast.error(serverMessage);
        }

        return Promise.reject(error);
    }
);

// Güvenli Çıkış (Logout) ve Yönlendirme Fonksiyonu
const forceLogout = () => {
    const wasAuthenticated = localStorage.getItem("isAuthenticated") === "true";
    localStorage.removeItem("isAuthenticated");
    localStorage.removeItem("isAdmin");
    localStorage.removeItem("userName");
    if (wasAuthenticated && window.location.pathname !== "/login") {
        window.location.href = "/login";
    }
}