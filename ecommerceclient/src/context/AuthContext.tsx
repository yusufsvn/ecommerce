import { createContext, useState, useContext, useEffect, type ReactNode } from 'react';
import { api } from '../Interceptors/httpinterceptor';

interface AuthContextType {
    isAuthenticated: boolean;
    isAdmin: boolean;
    userName: string | null;
    login: (name: string, asAdmin?: boolean) => void;
    logout: () => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider = ({ children }: { children: ReactNode }) => {
    const [isAuthenticated, setIsAuthenticated] = useState<boolean>(() => {
        return localStorage.getItem('isAuthenticated') === 'true';
    });

    const [isAdmin, setIsAdmin] = useState<boolean>(() => {
        return localStorage.getItem('isAdmin') === 'true';
    });

    const [userName, setUserName] = useState<string | null>(() => {
        return localStorage.getItem('userName');
    });

    const login = (name: string, asAdmin: boolean = false) => {
        setIsAuthenticated(true);
        setIsAdmin(asAdmin);
        setUserName(name);

        localStorage.setItem('isAuthenticated', 'true');
        localStorage.setItem('isAdmin', String(asAdmin));
        localStorage.setItem('userName', name);
    };

    const logout = () => {
        setIsAuthenticated(false);
        setIsAdmin(false);
        setUserName(null);

        localStorage.removeItem('isAuthenticated');
        localStorage.removeItem('isAdmin');
        localStorage.removeItem('userName');
    };

    useEffect(() => {
        // Uygulama yüklendiğinde arka plana sorarak kimliği doğrula
        const verifyUser = async () => {
            try {
                const response = await api.get('/user/me');

                if (response.status === 200) {
                    setIsAuthenticated(true);
                }
            } catch (error: any) {
                if (error.response?.status === 401 || error.response?.status === 403) {
                    console.error("Token geçersiz veya yetkisiz erişim. Çıkış yapılıyor.");
                    logout();
                } else {
                    // Diğer hatalarda (sunucu kapalı, ağ hatası vb.) logout yapma.
                    console.error("Backend'e ulaşılamadı veya bir hata oluştu:", error);
                }
            }
        };

        // Eğer kullanıcı zaten giriş yapmış gibi görünüyorsa teyit et
        if (localStorage.getItem('isAuthenticated') === 'true') {
            verifyUser();
        }
    }, []);

    return (
        <AuthContext.Provider value={{ isAuthenticated, isAdmin, userName, login, logout }}>
            {children}
        </AuthContext.Provider>
    );
};

export const useAuth = () => {
    const context = useContext(AuthContext);
    if (context === undefined) {
        throw new Error('useAuth must be used within an AuthProvider');
    }
    return context;
};
