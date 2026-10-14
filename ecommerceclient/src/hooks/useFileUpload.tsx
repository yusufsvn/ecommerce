import { useState } from "react";
import { fileService } from "../Services/FileService";

function useFileUpload() {

    const [files, setFiles] = useState<UploadingFile[]>([]);
    const [isChangedThumb, setIsChangedThumb] = useState<boolean>(false);

    const handleFileSelect = async (event: React.ChangeEvent<HTMLInputElement>) => {
        if (event.target.files) {
            const files = Array.from(event.target.files);

            // dosya objesi
            const newFiles = files.map(file => ({
                id: Math.random().toString(36).substring(2, 7),
                file,
                preview: URL.createObjectURL(file),
                progress: 0,
                status: 'uploading' as const
            }));

            setFiles(prev => [...prev, ...newFiles]);


            // İşlemin uzunluğuna göre progress bar
            for (const fileObj of newFiles) {
                try {
                    // işlem başladığını belirten yüzde 30 luk kısım
                    setFiles(prev => prev.map(file => {
                        if (file.id === fileObj.id) {
                            return { ...file, progress: 30 };
                        }
                        return file;
                    }));

                    // sıkıştırma işlemi
                    const compressedFile = await fileService.compressFile(fileObj.file);

                    // sıkıştırma işlemi bittikten sonra progress bar
                    setFiles(prev => prev.map(file => {
                        if (file.id === fileObj.id) {
                            return {
                                ...file,
                                file: compressedFile,
                                progress: 100,
                                status: 'completed'
                            };
                        }
                        return file;
                    }));
                } catch (error) {
                    console.error("Compression failed for", fileObj.file.name, error);
                    setFiles(prev => prev.map(file => {
                        if (file.id === fileObj.id) {
                            return { ...file, status: 'error', progress: 0 };
                        }
                        return file;
                    }));
                }
            }
        }
    };

    const handleRemoveFile = (id: string) => {
        setFiles(prev => {
            const index = prev.findIndex(f => f.id === id);
            if (index === 0 && isChangedThumb) {
                setIsChangedThumb(false);
            }
            return prev.filter(f => f.id !== id);
        });
    };


    const handleThumbnailSelect = async (event: React.ChangeEvent<HTMLInputElement>) => {
        if (event.target.files) {
            // seçilen fotoğrafların listesi
            const files = Array.from(event.target.files);

            // sadece 1 tanesini alsın
            const file = files[0];
            // dosya objesi
            const newFile = {
                id: Math.random().toString(36).substring(2, 7),
                file,
                preview: URL.createObjectURL(file),
                progress: 0,
                status: 'uploading' as const
            };

            // listenin en başına eklesin
            setFiles(prev => {
                if (isChangedThumb) {
                    return [newFile, ...prev.slice(1)];
                }
                return [newFile, ...prev];
            });
            setIsChangedThumb(true);

            try {
                // işlem başladığını belirten yüzde 30 luk kısım
                setFiles(prev => prev.map(file => {
                    if (file.id === newFile.id) {
                        return { ...file, progress: 30 };
                    }
                    return file;
                }));

                // sıkıştırma işlemi
                const compressedFile = await fileService.compressFile(newFile.file);

                // sıkıştırma işlemi bittikten sonra progress bar
                setFiles(prev => prev.map(file => {
                    if (file.id === newFile.id) {
                        return {
                            ...file,
                            file: compressedFile,
                            progress: 100,
                            status: 'completed'
                        };
                    }
                    return file;
                }));
            } catch (error) {
                console.error("Compression failed for", newFile.file.name, error);
                setFiles(prev => prev.map(file => {
                    if (file.id === newFile.id) {
                        return { ...file, status: 'error', progress: 0 };
                    }
                    return file;
                }));
            }

        }
    };

    const handleRemoveThumbnail = (id: string) => {
        setFiles(prev => prev.filter(f => f.id !== id));
    };

    const resetFiles = () => {
        setFiles([]);
        setIsChangedThumb(false);
    };

    return {
        files, handleFileSelect, handleRemoveFile, handleThumbnailSelect, handleRemoveThumbnail, isChangedThumb, resetFiles
    }
}

export default useFileUpload;

export interface UploadingFile {
    id: string;
    file: File;
    preview: string;
    progress: number;
    status: 'uploading' | 'completed' | 'error';
}