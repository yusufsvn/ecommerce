import { Box, Container, Typography, Button, Paper } from '@mui/material';
import { Grid } from '@mui/material';
import LockOutlinedIcon from '@mui/icons-material/LockOutlined';
import VerifiedUserOutlinedIcon from '@mui/icons-material/VerifiedUserOutlined';
import CreditCardOutlinedIcon from '@mui/icons-material/CreditCardOutlined';
import ReplayOutlinedIcon from '@mui/icons-material/ReplayOutlined';
import ArrowForwardIcon from '@mui/icons-material/ArrowForward';
import { useNavigate } from 'react-router';
import securePaymentImg from '../assets/secure_payment.jpg';

export default function SecurePaymentSection() {
  const navigate = useNavigate();

  const securityFeatures = [
    {
      icon: <LockOutlinedIcon sx={{ color: 'primary.main', fontSize: 24 }} />,
      title: '256-Bit SSL Sertifikası',
      desc: 'Banka seviyesinde veri şifreleme protokolü',
    },
    {
      icon: <VerifiedUserOutlinedIcon sx={{ color: 'primary.main', fontSize: 24 }} />,
      title: '3D Secure Doğrulama',
      desc: 'Tek kullanımlık SMS şifresiyle onaylı ödeme',
    },
    {
      icon: <CreditCardOutlinedIcon sx={{ color: 'primary.main', fontSize: 24 }} />,
      title: 'Tüm Kartlara Taksit',
      desc: 'Anlaşmalı kredi kartlarına vade farksız seçenekler',
    },
    {
      icon: <ReplayOutlinedIcon sx={{ color: 'primary.main', fontSize: 24 }} />,
      title: '14 Gün Kolay İade',
      desc: 'Koşulsuz ve kesintisiz iade garantisi',
    },
  ];

  return (
    <Box sx={{ bgcolor: 'background.default', py: { xs: 6, md: 10 } }}>
      <Container maxWidth="lg" sx={{ px: { xs: 2, md: 6 } }}>
        <Paper
          elevation={0}
          sx={{
            borderRadius: 5,
            border: '1px solid',
            borderColor: 'divider',
            bgcolor: 'background.paper',
            overflow: 'hidden',
            p: { xs: 3, sm: 4, md: 6 },
            boxShadow: '0 10px 30px -10px rgba(30, 58, 138, 0.08)',
          }}
        >
          <Grid container spacing={{ xs: 4, md: 6 }} alignItems="center">
            {/* Sol Taraf: Güvenli Ödeme Fotoğrafı */}
            <Grid size={{ xs: 12, md: 6 }}>
              <Box
                sx={{
                  position: 'relative',
                  width: '100%',
                  borderRadius: 4,
                  overflow: 'hidden',
                  boxShadow: '0 20px 40px -15px rgba(30, 58, 138, 0.22)',
                  '&:hover img': {
                    transform: 'scale(1.04)',
                  },
                }}
              >
                <Box
                  component="img"
                  src={securePaymentImg}
                  alt="Güvenli Ödeme Altyapısı"
                  loading="lazy"
                  sx={{
                    width: '100%',
                    height: { xs: 280, sm: 380, md: 440 },
                    objectFit: 'cover',
                    display: 'block',
                    transition: 'transform 0.6s cubic-bezier(0.25, 1, 0.5, 1)',
                  }}
                />

                {/* Fotoğraf Üzerinde Yüzen Güvenlik Rozeti */}
                <Box
                  sx={{
                    position: 'absolute',
                    bottom: 20,
                    left: 20,
                    right: { xs: 20, sm: 'auto' },
                    bgcolor: 'rgba(255, 255, 255, 0.92)',
                    backdropFilter: 'blur(12px)',
                    borderRadius: 3,
                    p: 2,
                    display: 'flex',
                    alignItems: 'center',
                    gap: 1.5,
                    boxShadow: '0 8px 24px rgba(0, 0, 0, 0.15)',
                    border: '1px solid rgba(255, 255, 255, 0.6)',
                  }}
                >
                  <Box
                    sx={{
                      width: 42,
                      height: 42,
                      borderRadius: 2.5,
                      bgcolor: 'primary.main',
                      color: 'white',
                      display: 'flex',
                      alignItems: 'center',
                      justifyContent: 'center',
                      flexShrink: 0,
                    }}
                  >
                    <LockOutlinedIcon sx={{ fontSize: 22 }} />
                  </Box>
                  <Box>
                    <Typography variant="subtitle2" sx={{ fontWeight: 800, color: 'primary.dark', lineHeight: 1.2 }}>
                      256-Bit SSL Koruma
                    </Typography>
                    <Typography variant="caption" sx={{ color: 'text.secondary', fontWeight: 500 }}>
                      Kart bilgileriniz asla saklanmaz
                    </Typography>
                  </Box>
                </Box>
              </Box>
            </Grid>

            {/* Sağ Taraf: Slogan Başlık, Açıklama ve Güvenlik Maddeleri */}
            <Grid size={{ xs: 12, md: 6 }}>
              <Box>
                {/* Üst Vurgu Rozeti */}
                <Box
                  sx={{
                    display: 'inline-flex',
                    alignItems: 'center',
                    gap: 1,
                    px: 1.75,
                    py: 0.6,
                    borderRadius: 32,
                    bgcolor: 'rgba(249, 115, 22, 0.1)',
                    color: 'secondary.main',
                    mb: 2,
                  }}
                >
                  <VerifiedUserOutlinedIcon sx={{ fontSize: 18 }} />
                  <Typography variant="caption" sx={{ fontWeight: 800, letterSpacing: '1px', textTransform: 'uppercase' }}>
                    Güvenli Ödeme Altyapısı
                  </Typography>
                </Box>

                {/* Slogan Başlık */}
                <Typography
                  variant="h3"
                  component="h2"
                  sx={{
                    fontWeight: 800,
                    color: 'primary.main',
                    fontSize: { xs: '1.75rem', sm: '2.15rem', md: '2.5rem' },
                    lineHeight: 1.2,
                    letterSpacing: '-0.5px',
                    mb: 2,
                  }}
                >
                  Ödemeleriniz %100 Güvende, Alışverişiniz Keyifle.
                </Typography>

                {/* Kısa Altyazı */}
                <Typography
                  variant="body1"
                  sx={{
                    color: 'text.secondary',
                    fontSize: { xs: '0.95rem', md: '1.05rem' },
                    lineHeight: 1.7,
                    mb: 3.5,
                  }}
                >
                  Sitemiz üzerinden gerçekleştirdiğiniz tüm siparişler, uluslararası güvenlik standartlarına sahip 256-bit SSL ve 3D Secure sistemleri ile uçtan uca korunmaktadır.
                </Typography>

                {/* Özellik Maddeleri Grid */}
                <Grid container spacing={2} sx={{ mb: 4 }}>
                  {securityFeatures.map((item, index) => (
                    <Grid size={{ xs: 12, sm: 6 }} key={index}>
                      <Box
                        sx={{
                          display: 'flex',
                          alignItems: 'flex-start',
                          gap: 1.5,
                          p: 1.5,
                          borderRadius: 2.5,
                          bgcolor: 'background.default',
                          border: '1px solid',
                          borderColor: 'divider',
                        }}
                      >
                        <Box
                          sx={{
                            width: 36,
                            height: 36,
                            borderRadius: 2,
                            bgcolor: 'background.paper',
                            display: 'flex',
                            alignItems: 'center',
                            justifyContent: 'center',
                            flexShrink: 0,
                            boxShadow: '0 2px 6px rgba(0,0,0,0.05)',
                          }}
                        >
                          {item.icon}
                        </Box>
                        <Box>
                          <Typography variant="body2" sx={{ fontWeight: 700, color: 'text.primary', lineHeight: 1.3 }}>
                            {item.title}
                          </Typography>
                          <Typography variant="caption" sx={{ color: 'text.secondary', display: 'block', mt: 0.3 }}>
                            {item.desc}
                          </Typography>
                        </Box>
                      </Box>
                    </Grid>
                  ))}
                </Grid>

                {/* Aksiyon Butonu & Desteklenen Kart Logoları */}
                <Box sx={{ display: 'flex', flexWrap: 'wrap', alignItems: 'center', gap: 3 }}>
                  <Button
                    variant="contained"
                    color="secondary"
                    size="large"
                    endIcon={<ArrowForwardIcon />}
                    onClick={() => navigate('/products')}
                    sx={{
                      px: 4,
                      py: 1.4,
                      fontWeight: 700,
                      fontSize: '1rem',
                      boxShadow: '0 6px 18px rgba(249, 115, 22, 0.35)',
                      '&:hover': {
                        boxShadow: '0 8px 25px rgba(249, 115, 22, 0.45)',
                      },
                    }}
                  >
                    Alışverişe Başla
                  </Button>

                  <Box sx={{ display: 'flex', alignItems: 'center', gap: 1.5, opacity: 0.75 }}>
                    <Typography variant="caption" sx={{ fontWeight: 600, color: 'text.secondary' }}>
                      Kabul Edilen Kartlar:
                    </Typography>
                    <Box sx={{ display: 'flex', gap: 1 }}>
                      <Box
                        sx={{
                          px: 1.2,
                          py: 0.4,
                          borderRadius: 1.5,
                          border: '1px solid',
                          borderColor: 'divider',
                          bgcolor: 'background.paper',
                          fontWeight: 800,
                          fontSize: '0.75rem',
                          color: '#1A1F71',
                        }}
                      >
                        VISA
                      </Box>
                      <Box
                        sx={{
                          px: 1.2,
                          py: 0.4,
                          borderRadius: 1.5,
                          border: '1px solid',
                          borderColor: 'divider',
                          bgcolor: 'background.paper',
                          fontWeight: 800,
                          fontSize: '0.75rem',
                          color: '#EB001B',
                        }}
                      >
                        Mastercard
                      </Box>
                      <Box
                        sx={{
                          px: 1.2,
                          py: 0.4,
                          borderRadius: 1.5,
                          border: '1px solid',
                          borderColor: 'divider',
                          bgcolor: 'background.paper',
                          fontWeight: 800,
                          fontSize: '0.75rem',
                          color: '#0057A0',
                        }}
                      >
                        Troy
                      </Box>
                    </Box>
                  </Box>
                </Box>
              </Box>
            </Grid>
          </Grid>
        </Paper>
      </Container>
    </Box>
  );
}
