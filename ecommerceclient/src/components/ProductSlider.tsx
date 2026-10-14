import { useState } from 'react';
import { Box, Paper, Button, MobileStepper } from '@mui/material';
import { useTheme } from '@mui/material/styles';
import KeyboardArrowLeft from '@mui/icons-material/KeyboardArrowLeft';
import KeyboardArrowRight from '@mui/icons-material/KeyboardArrowRight';
import type { Image } from '../models/ProductResponse';

function ProductSlider({ images }: { images: Image[] }) {
    const theme = useTheme();
    const [activeStep, setActiveStep] = useState(0);
    const maxSteps = images.length;

    const handleNext = () => {
        setActiveStep((prevActiveStep) => prevActiveStep + 1);
    };

    const handleBack = () => {
        setActiveStep((prevActiveStep) => prevActiveStep - 1);
    };

    // Eğer resim yoksa boş bir kutu göster (Hata önleyici)
    if (!images || images.length === 0) {
        return <Box sx={{ p: 3, textAlign: 'center' }}>Görsel Yok</Box>;
    }

    console.log(images[0].image_id)

    return (
        <Box sx={{ maxWidth: '100%', flexGrow: 1 }}>
            {/* Görsel Alanı */}
            <Paper
                square
                elevation={0}
                sx={{
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'center',
                    height: 400, // Slider yüksekliği
                    bgcolor: 'background.default',
                    overflow: 'hidden',
                    position: 'relative' // Okları resmin üzerine koymak istersen bu gerekli
                }}
            >
                {/* Mevcut aktif resmi göster */}
                <Box
                    component="img"
                    sx={{
                        height: '100%',
                        maxWidth: '100%',
                        display: 'block',
                        objectFit: 'contain', // Resmi kırpmadan sığdırır
                    }}
                    src={images[activeStep].path}
                    alt={images[activeStep].image_id || `product-img-${activeStep}`}
                />
            </Paper>

            {/* Navigasyon (Noktalar ve Oklar) */}
            <MobileStepper
                variant="dots" // "text", "dots", veya "progress" olabilir
                steps={maxSteps}
                position="static"
                activeStep={activeStep}
                sx={{
                    maxWidth: '100%',
                    flexGrow: 1,
                    bgcolor: 'transparent' // Arkaplan rengini temizle
                }}
                nextButton={
                    <Button
                        size="small"
                        onClick={handleNext}
                        disabled={activeStep === maxSteps - 1}
                    >
                        İleri
                        {theme.direction === 'rtl' ? (
                            <KeyboardArrowLeft />
                        ) : (
                            <KeyboardArrowRight />
                        )}
                    </Button>
                }
                backButton={
                    <Button size="small" onClick={handleBack} disabled={activeStep === 0}>
                        {theme.direction === 'rtl' ? (
                            <KeyboardArrowRight />
                        ) : (
                            <KeyboardArrowLeft />
                        )}
                        Geri
                    </Button>
                }
            />
        </Box>
    );
};

export default ProductSlider;

