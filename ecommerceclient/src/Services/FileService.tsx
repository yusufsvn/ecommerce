import imageCompression from "browser-image-compression";

// E-Ticaret için İdeal Ayarlar (1 Core Sunucu Dostu)
const options = {
    maxSizeMB: 1,               // Hedef: 1MB altı
    maxWidthOrHeight: 1920,     // Full HD yeterli
    useWebWorker: true,         // Tarayıcıyı dondurmasın
    fileType: 'image/webp',     // Format dönüşümü
    initialQuality: 0.8,        // Kalite
};

export const fileService = {

    async compressFile(file: File): Promise<File> {
        try {
            console.log(`Orijinal Boyut: ${(file.size / 1024 / 1024).toFixed(2)} MB`);

            // AĞIR İŞLEM BURADA (Kullanıcının CPU'sunda) YAPILIYOR
            const compressedFile = await imageCompression(file, options);

            console.log(`Sıkıştırılmış Boyut: ${(compressedFile.size / 1024 / 1024).toFixed(2)} MB`);

            return compressedFile;

        } catch (error) {
            console.error("Sıkıştırma hatası:", error);
            throw error;
        }
    }



}
