import React, { useState } from 'react';
import {
  Container,
  Typography,
  Box,
  Button,
  IconButton,
  OutlinedInput,
  InputAdornment,
  Checkbox,
  FormControlLabel,
  Divider,
  Link,
  Paper,
  Grid,
} from '@mui/material';
import Visibility from '@mui/icons-material/Visibility';
import VisibilityOff from '@mui/icons-material/VisibilityOff';
import ArrowForwardIcon from '@mui/icons-material/ArrowForward';
import CheckCircleIcon from '@mui/icons-material/CheckCircle';
import { useNavigate } from 'react-router';
import { useAuth } from '../context/AuthContext';
import { userService } from '../Services/userService';
import { toast } from 'react-toastify';
import { GoogleLogin } from '@react-oauth/google';

const Login = () => {
  const navigate = useNavigate();
  const { login } = useAuth();
  const [formData, setFormData] = useState({
    email: '',
    password: '',
  });
  const [emailError, setEmailError] = useState('');
  const [showPassword, setShowPassword] = useState(false);
  const [rememberMe, setRememberMe] = useState(false);
  const [loading, setLoading] = useState(false);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setFormData((prev) => ({ ...prev, [name]: value }));
    if (name === 'email' && emailError) {
      setEmailError('');
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    const emailRegex = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/;
    if (!emailRegex.test(formData.email)) {
      setEmailError('Lütfen geçerli bir e-posta adresi giriniz.');
      return;
    } else {
      setEmailError('');
    }

    if (!formData.password) {
      toast.warning('Lütfen şifrenizi giriniz.');
      return;
    }

    setLoading(true);
    try {
      const data = await userService.loginUser({ mail: formData.email, password: formData.password });
      if (data) {
        login(data.name + ' ' + data.surname);
        toast.success('Hoş geldiniz, ' + data.name + '!');
        navigate('/');
      }
    } catch (error) {
      toast.error('Hatalı giriş! Lütfen bilgilerinizi kontrol edin.');
    } finally {
      setLoading(false);
    }
  };

  const handleSuccess = async (credentialResponse: any) => {
    try {
      const idToken = credentialResponse.credential;
      const data = await userService.loginWithGoogle(idToken);
      if (data) {
        login(data.name + ' ' + data.surname);
        toast.success('Hoş geldiniz, ' + data.name + '!');
        navigate('/');
      }
    } catch (error) {
      toast.error('Google ile giriş başarısız oldu. Lütfen tekrar deneyin.');
    }
  };

  const handleError = () => {
    toast.error('Google ile giriş başarısız oldu. Lütfen tekrar deneyin.');
  };

  return (
    <Box
      sx={{
        bgcolor: '#F8FAFC',
        py: { xs: 3, md: 6 },
        minHeight: 'calc(100vh - 140px)',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
      }}
    >
      <Container maxWidth="lg">
        <Paper
          elevation={0}
          sx={{
            borderRadius: { xs: 4, md: 6 },
            overflow: 'hidden',
            bgcolor: '#FFFFFF',
            boxShadow: '0 20px 45px -15px rgba(15, 23, 42, 0.08), 0 1px 3px rgba(15, 23, 42, 0.04)',
            border: '1px solid rgba(226, 232, 240, 0.8)',
            p: { xs: 2, md: 2.5 },
          }}
        >
          <Grid container spacing={{ xs: 3, md: 4 }} alignItems="stretch">
            {/* SOL SÜTUN - GÖRSEL VE VİZYON ALANI */}
            <Grid size={{ xs: 12, md: 6 }}>
              <Box
                sx={{
                  position: 'relative',
                  height: '100%',
                  minHeight: { xs: 360, md: 620 },
                  borderRadius: { xs: 3, md: 5 },
                  overflow: 'hidden',
                  backgroundImage: 'url(/login-showcase.jpg)',
                  backgroundSize: 'cover',
                  backgroundPosition: 'center',
                  display: 'flex',
                  flexDirection: 'column',
                  justifyContent: 'flex-end',
                  boxShadow: 'inset 0 0 0 1px rgba(255, 255, 255, 0.1)',
                }}
              >
                {/* Gradient Karartma Katmanı */}
                <Box
                  sx={{
                    position: 'absolute',
                    inset: 0,
                    background:
                      'linear-gradient(180deg, rgba(15, 23, 42, 0.05) 0%, rgba(15, 23, 42, 0.3) 35%, rgba(11, 20, 44, 0.9) 100%)',
                    zIndex: 1,
                  }}
                />

                {/* İçerik */}
                <Box
                  sx={{
                    position: 'relative',
                    zIndex: 2,
                    p: { xs: 3, sm: 4, md: 4.5 },
                  }}
                >
                  {/* Üst Rozet / Pill */}
                  <Box
                    sx={{
                      display: 'inline-flex',
                      alignItems: 'center',
                      gap: 1.2,
                      px: 2,
                      py: 0.7,
                      borderRadius: '999px',
                      bgcolor: 'rgba(255, 255, 255, 0.12)',
                      backdropFilter: 'blur(12px)',
                      border: '1px solid rgba(255, 255, 255, 0.25)',
                      mb: 2.5,
                    }}
                  >
                    <Box
                      sx={{
                        width: 8,
                        height: 8,
                        borderRadius: '50%',
                        bgcolor: '#F97316',
                        boxShadow: '0 0 8px #F97316',
                      }}
                    />
                    <Typography
                      sx={{
                        fontSize: '0.78rem',
                        fontWeight: 700,
                        color: '#FFFFFF',
                        fontFamily: '"Inter", sans-serif',
                        letterSpacing: '0.2px',
                      }}
                    >
                      Ayrıcalıklı Alışveriş Deneyimi
                    </Typography>
                  </Box>

                  {/* Başlık */}
                  <Typography
                    variant="h3"
                    component="h2"
                    sx={{
                      fontFamily: '"Outfit", sans-serif',
                      fontWeight: 800,
                      color: '#FFFFFF',
                      fontSize: { xs: '1.6rem', sm: '1.9rem', md: '2.15rem' },
                      lineHeight: 1.25,
                      letterSpacing: '-0.5px',
                      mb: 1.5,
                    }}
                  >
                    Teknoloji ve Estetiğin Buluşma Noktası
                  </Typography>

                  {/* Açıklama */}
                  <Typography
                    variant="body2"
                    sx={{
                      fontFamily: '"Inter", sans-serif',
                      color: 'rgba(255, 255, 255, 0.85)',
                      fontSize: { xs: '0.82rem', sm: '0.88rem' },
                      lineHeight: 1.55,
                      mb: 3,
                      maxWidth: '96%',
                    }}
                  >
                    Binlerce seçkin ürün, ayrıcalıklı fırsatlar ve güvenli alışveriş deneyimi Ecommerce Store'da sizi bekliyor.
                  </Typography>

                  {/* Özellik Rozetleri */}
                  <Box sx={{ display: 'flex', alignItems: 'center', gap: 3, flexWrap: 'wrap' }}>
                    <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                      <CheckCircleIcon sx={{ fontSize: 19, color: '#F97316' }} />
                      <Typography
                        sx={{
                          color: '#FFFFFF',
                          fontWeight: 700,
                          fontSize: '0.84rem',
                          fontFamily: '"Inter", sans-serif',
                        }}
                      >
                        Orijinal Ürünler
                      </Typography>
                    </Box>

                    <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                      <CheckCircleIcon sx={{ fontSize: 19, color: '#F97316' }} />
                      <Typography
                        sx={{
                          color: '#FFFFFF',
                          fontWeight: 700,
                          fontSize: '0.84rem',
                          fontFamily: '"Inter", sans-serif',
                        }}
                      >
                        Güvenli Ödeme
                      </Typography>
                    </Box>
                  </Box>
                </Box>
              </Box>
            </Grid>

            {/* SAĞ SÜTUN - GİRİŞ FORMU */}
            <Grid
              size={{ xs: 12, md: 6 }}
              sx={{
                display: 'flex',
                flexDirection: 'column',
                justifyContent: 'center',
                px: { xs: 2, sm: 4, md: 5 },
                py: { xs: 3, md: 4 },
              }}
            >
              <Box sx={{ maxWidth: 440, mx: 'auto', width: '100%' }}>
                {/* Form Başlıkları */}
                <Box sx={{ mb: 3.5 }}>
                  <Typography
                    component="h1"
                    variant="h4"
                    sx={{
                      fontFamily: '"Outfit", sans-serif',
                      fontWeight: 800,
                      color: '#1E3A8A',
                      fontSize: { xs: '2rem', md: '2.35rem' },
                      letterSpacing: '-0.6px',
                      mb: 1,
                    }}
                  >
                    Giriş Yap
                  </Typography>
                  <Typography
                    variant="body2"
                    sx={{
                      color: '#64748B',
                      fontFamily: '"Inter", sans-serif',
                      fontSize: '0.92rem',
                    }}
                  >
                    Hesabınıza erişmek için lütfen bilgilerinizi girin.
                  </Typography>
                </Box>

                {/* Form */}
                <Box component="form" onSubmit={handleSubmit} noValidate>
                  {/* E-posta Alanı */}
                  <Box sx={{ mb: 2.5 }}>
                    <Typography
                      component="label"
                      htmlFor="email"
                      sx={{
                        display: 'block',
                        fontSize: '0.72rem',
                        fontWeight: 800,
                        letterSpacing: '0.6px',
                        color: '#1E293B',
                        fontFamily: '"Inter", sans-serif',
                        mb: 0.8,
                      }}
                    >
                      E-POSTA ADRESİ <span style={{ color: '#F97316' }}>*</span>
                    </Typography>
                    <OutlinedInput
                      id="email"
                      name="email"
                      type="email"
                      autoComplete="email"
                      autoFocus
                      fullWidth
                      value={formData.email}
                      onChange={handleChange}
                      placeholder="ornek@hesap.com"
                      error={!!emailError}
                      sx={{
                        borderRadius: '999px',
                        bgcolor: '#FFFFFF',
                        fontSize: '0.92rem',
                        fontFamily: '"Inter", sans-serif',
                        '& fieldset': {
                          borderColor: emailError ? '#EF4444' : '#E2E8F0',
                        },
                        '&:hover fieldset': {
                          borderColor: emailError ? '#EF4444' : '#94A3B8',
                        },
                        '&.Mui-focused fieldset': {
                          borderColor: emailError ? '#EF4444' : '#1E3A8A',
                          borderWidth: '1.5px',
                        },
                        '& input': {
                          px: 3,
                          py: 1.4,
                        },
                      }}
                    />
                    {emailError && (
                      <Typography
                        variant="caption"
                        sx={{ color: '#EF4444', mt: 0.5, display: 'block', pl: 1.5, fontWeight: 500 }}
                      >
                        {emailError}
                      </Typography>
                    )}
                  </Box>

                  {/* Şifre Alanı */}
                  <Box sx={{ mb: 2 }}>
                    <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 0.8 }}>
                      <Typography
                        component="label"
                        htmlFor="password"
                        sx={{
                          fontSize: '0.72rem',
                          fontWeight: 800,
                          letterSpacing: '0.6px',
                          color: '#1E293B',
                          fontFamily: '"Inter", sans-serif',
                        }}
                      >
                        ŞİFRE <span style={{ color: '#F97316' }}>*</span>
                      </Typography>
                      <Link
                        component="button"
                        type="button"
                        onClick={() => toast.info('Şifre sıfırlama işlemi için lütfen destek ekibiyle iletişime geçin.')}
                        sx={{
                          fontSize: '0.8rem',
                          fontWeight: 700,
                          color: '#EA580C',
                          textDecoration: 'none',
                          fontFamily: '"Inter", sans-serif',
                          cursor: 'pointer',
                          '&:hover': { textDecoration: 'underline', color: '#C2410C' },
                        }}
                      >
                        Şifremi Unuttum?
                      </Link>
                    </Box>
                    <OutlinedInput
                      id="password"
                      name="password"
                      type={showPassword ? 'text' : 'password'}
                      autoComplete="current-password"
                      fullWidth
                      value={formData.password}
                      onChange={handleChange}
                      placeholder="••••••••"
                      endAdornment={
                        <InputAdornment position="end" sx={{ mr: 1 }}>
                          <IconButton
                            aria-label="şifre görünürlüğünü değiştir"
                            onClick={() => setShowPassword(!showPassword)}
                            edge="end"
                            size="small"
                            sx={{ color: '#94A3B8', '&:hover': { color: '#1E3A8A' } }}
                          >
                            {showPassword ? <VisibilityOff fontSize="small" /> : <Visibility fontSize="small" />}
                          </IconButton>
                        </InputAdornment>
                      }
                      sx={{
                        borderRadius: '999px',
                        bgcolor: '#FFFFFF',
                        fontSize: '0.92rem',
                        fontFamily: '"Inter", sans-serif',
                        '& fieldset': {
                          borderColor: '#E2E8F0',
                        },
                        '&:hover fieldset': {
                          borderColor: '#94A3B8',
                        },
                        '&.Mui-focused fieldset': {
                          borderColor: '#1E3A8A',
                          borderWidth: '1.5px',
                        },
                        '& input': {
                          px: 3,
                          py: 1.4,
                        },
                      }}
                    />
                  </Box>



                  {/* Giriş Yap Butonu */}
                  <Button
                    type="submit"
                    fullWidth
                    disabled={loading}
                    endIcon={<ArrowForwardIcon sx={{ fontSize: '1.15rem !important' }} />}
                    sx={{
                      bgcolor: '#1E3A8A',
                      color: '#FFFFFF',
                      borderRadius: '999px',
                      py: 1.5,
                      fontSize: '0.9rem',
                      fontWeight: 800,
                      letterSpacing: '0.5px',
                      fontFamily: '"Inter", sans-serif',
                      boxShadow: '0 8px 20px -4px rgba(30, 58, 138, 0.35)',
                      transition: 'all 0.25s ease',
                      '&:hover': {
                        bgcolor: '#172554',
                        transform: 'translateY(-1px)',
                        boxShadow: '0 12px 24px -4px rgba(30, 58, 138, 0.45)',
                      },
                    }}
                  >
                    {loading ? 'GİRİŞ YAPILIYOR...' : 'GİRİŞ YAP'}
                  </Button>

                  {/* VEYA Ayırıcı */}
                  <Box sx={{ display: 'flex', alignItems: 'center', my: 2.5 }}>
                    <Divider sx={{ flex: 1, borderColor: '#E2E8F0' }} />
                    <Typography
                      sx={{
                        px: 2,
                        fontSize: '0.72rem',
                        fontWeight: 800,
                        color: '#94A3B8',
                        letterSpacing: '0.8px',
                        fontFamily: '"Inter", sans-serif',
                      }}
                    >
                      VEYA
                    </Typography>
                    <Divider sx={{ flex: 1, borderColor: '#E2E8F0' }} />
                  </Box>

                  {/* Google ile Giriş Butonu */}
                  <Box
                    sx={{
                      display: 'flex',
                      justifyContent: 'center',
                      width: '100%',
                      '& > div': { width: '100% !important' },
                      '& iframe': { margin: '0 auto !important' },
                    }}
                  >
                    <GoogleLogin
                      onSuccess={handleSuccess}
                      onError={handleError}
                      shape="pill"
                      size="large"
                      theme="outline"
                      text="signin_with"
                      width="380"
                    />
                  </Box>

                  {/* Kayıt Ol Linki */}
                  <Box sx={{ display: 'flex', justifyContent: 'center', mt: 3.5 }}>
                    <Typography variant="body2" sx={{ color: '#64748B', fontSize: '0.88rem', fontWeight: 500 }}>
                      Hesabın yok mu?{' '}
                      <Link
                        component="button"
                        type="button"
                        onClick={() => navigate('/register')}
                        sx={{
                          color: '#1E3A8A',
                          fontWeight: 800,
                          fontSize: '0.88rem',
                          textDecoration: 'none',
                          cursor: 'pointer',
                          fontFamily: '"Inter", sans-serif',
                          '&:hover': { textDecoration: 'underline', color: '#172554' },
                        }}
                      >
                        Hemen Kayıt Ol
                      </Link>
                    </Typography>
                  </Box>
                </Box>
              </Box>
            </Grid>
          </Grid>
        </Paper>
      </Container>
    </Box>
  );
};

export default Login;
