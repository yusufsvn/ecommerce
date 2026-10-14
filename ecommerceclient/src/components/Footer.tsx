import { Box, Container, Stack, Typography, Link as MuiLink, IconButton, Divider } from '@mui/material';
import ShareIcon from '@mui/icons-material/Share';
import GroupsIcon from '@mui/icons-material/Groups';
import LanguageIcon from '@mui/icons-material/Language';
import LocationOnOutlinedIcon from '@mui/icons-material/LocationOnOutlined';
import PhoneOutlinedIcon from '@mui/icons-material/PhoneOutlined';
import EmailOutlinedIcon from '@mui/icons-material/EmailOutlined';
import { Link as RouterLink } from 'react-router';

function Footer() {
  return (
    <Box
      component="footer"
      sx={{
        bgcolor: '#1E2433', // Fotoğraftaki zarif koyu lacivert/arduvaz zemin
        color: '#FFFFFF',
        pt: { xs: 6, md: 7 },
        pb: 4,
        mt: 'auto',
        borderTop: '1px solid #2A3447',
      }}
    >
      <Container maxWidth="lg">
        {/* Üst Kısım: 3 Kolonlu Yapı */}
        <Stack
          direction={{ xs: 'column', md: 'row' }}
          spacing={{ xs: 4, md: 6 }}
          justifyContent="space-between"
        >
          {/* 1. KOLON: Logo, Açıklama ve Petrol Mavisi Sosyal Butonlar */}
          <Box sx={{ flex: 1.2 }}>
            <Typography
              variant="h6"
              sx={{
                fontFamily: '"Outfit", sans-serif',
                fontWeight: 800,
                color: '#FFFFFF',
                letterSpacing: '-0.5px',
                fontSize: '1.25rem',
                mb: 2,
              }}
            >
              Ecommerce Store
            </Typography>

            <Typography
              variant="body2"
              sx={{
                color: '#94A3B8',
                mb: 3,
                maxWidth: '340px',
                lineHeight: 1.6,
                fontSize: '0.875rem',
              }}
            >
              Kalite ve estetiği bir araya getirerek yaşamınıza değer katan seçkin ürünler sunuyoruz.
              Güvenli ve keyifli alışverişin adresi.
            </Typography>

            {/* Fotoğraftaki Yuvarlak Petrol Mavisi Düğmeler */}
            <Stack direction="row" spacing={1.5}>
              <IconButton
                size="small"
                sx={{
                  color: '#FFFFFF',
                  bgcolor: '#0E7490',
                  width: 36,
                  height: 36,
                  transition: 'all 0.2s ease',
                  '&:hover': { bgcolor: '#0891B2', transform: 'translateY(-2px)' },
                }}
                aria-label="Paylaş"
              >
                <ShareIcon sx={{ fontSize: 18 }} />
              </IconButton>
              <IconButton
                size="small"
                sx={{
                  color: '#FFFFFF',
                  bgcolor: '#0E7490',
                  width: 36,
                  height: 36,
                  transition: 'all 0.2s ease',
                  '&:hover': { bgcolor: '#0891B2', transform: 'translateY(-2px)' },
                }}
                aria-label="Topluluk"
              >
                <GroupsIcon sx={{ fontSize: 18 }} />
              </IconButton>
              <IconButton
                size="small"
                sx={{
                  color: '#FFFFFF',
                  bgcolor: '#0E7490',
                  width: 36,
                  height: 36,
                  transition: 'all 0.2s ease',
                  '&:hover': { bgcolor: '#0891B2', transform: 'translateY(-2px)' },
                }}
                aria-label="Web Sitesi"
              >
                <LanguageIcon sx={{ fontSize: 18 }} />
              </IconButton>
            </Stack>
          </Box>

          {/* 2. KOLON: Kurumsal Bağlantılar (Beyaz Başlık) */}
          <Box sx={{ flex: 0.8 }}>
            <Typography
              variant="subtitle1"
              sx={{
                fontFamily: '"Outfit", sans-serif',
                fontWeight: 700,
                color: '#FFFFFF',
                fontSize: '1rem',
                mb: 2,
              }}
            >
              Kurumsal
            </Typography>
            <Stack spacing={1.2}>
              <MuiLink
                component={RouterLink}
                to="/"
                underline="hover"
                sx={{ color: '#94A3B8', fontSize: '0.875rem', '&:hover': { color: '#FFFFFF' } }}
              >
                Ana Sayfa
              </MuiLink>
              <MuiLink
                component={RouterLink}
                to="/hakkimizda"
                underline="hover"
                sx={{ color: '#94A3B8', fontSize: '0.875rem', '&:hover': { color: '#FFFFFF' } }}
              >
                Hakkımızda
              </MuiLink>
              <MuiLink
                component={RouterLink}
                to="/products"
                underline="hover"
                sx={{ color: '#94A3B8', fontSize: '0.875rem', '&:hover': { color: '#FFFFFF' } }}
              >
                Ürünler
              </MuiLink>
              <MuiLink
                component={RouterLink}
                to="/iletisim"
                underline="hover"
                sx={{ color: '#94A3B8', fontSize: '0.875rem', '&:hover': { color: '#FFFFFF' } }}
              >
                İletişim
              </MuiLink>
              <MuiLink
                component={RouterLink}
                to="/cart"
                underline="hover"
                sx={{ color: '#94A3B8', fontSize: '0.875rem', '&:hover': { color: '#FFFFFF' } }}
              >
                Sepetim
              </MuiLink>
            </Stack>
          </Box>

          {/* 3. KOLON: İletişim Bilgileri (Beyaz Başlık) */}
          <Box sx={{ flex: 1 }}>
            <Typography
              variant="subtitle1"
              sx={{
                fontFamily: '"Outfit", sans-serif',
                fontWeight: 700,
                color: '#FFFFFF',
                fontSize: '1rem',
                mb: 2,
              }}
            >
              Bize Ulaşın
            </Typography>
            <Stack spacing={2}>
              <Stack direction="row" spacing={1.5} alignItems="flex-start">
                <LocationOnOutlinedIcon sx={{ color: '#94A3B8', mt: 0.25, fontSize: '1.25rem' }} />
                <Typography variant="body2" sx={{ color: '#94A3B8', lineHeight: 1.5, fontSize: '0.875rem' }}>
                  Merkez Plaza, Kat: 3, No: 42
                  <br />
                  Nilüfer / Bursa
                </Typography>
              </Stack>

              <Stack direction="row" spacing={1.5} alignItems="center">
                <PhoneOutlinedIcon sx={{ color: '#94A3B8', fontSize: '1.25rem' }} />
                <Typography variant="body2" sx={{ color: '#94A3B8', fontSize: '0.875rem' }}>
                  +90 (224) 123 45 67
                </Typography>
              </Stack>

              <Stack direction="row" spacing={1.5} alignItems="center">
                <EmailOutlinedIcon sx={{ color: '#94A3B8', fontSize: '1.25rem' }} />
                <Typography variant="body2" sx={{ color: '#94A3B8', fontSize: '0.875rem' }}>
                  info@sirketiniz.com
                </Typography>
              </Stack>
            </Stack>
          </Box>
        </Stack>

        <Divider sx={{ my: 4, borderColor: 'rgba(255, 255, 255, 0.08)' }} />

        {/* Alt Kısım: Telif Hakkı & Hukuki Bağlantılar */}
        <Stack
          direction={{ xs: 'column', sm: 'row' }}
          justifyContent="space-between"
          alignItems="center"
          spacing={2}
        >
          <Typography variant="body2" sx={{ color: '#64748B', fontSize: '0.85rem' }}>
            © {new Date().getFullYear()} Şirket Adı. Tüm hakları saklıdır.
          </Typography>
          <Stack direction="row" spacing={3}>
            <MuiLink
              href="#"
              underline="hover"
              sx={{ color: '#64748B', fontSize: '0.85rem', '&:hover': { color: '#94A3B8' } }}
            >
              Gizlilik Politikası
            </MuiLink>
            <MuiLink
              href="#"
              underline="hover"
              sx={{ color: '#64748B', fontSize: '0.85rem', '&:hover': { color: '#94A3B8' } }}
            >
              Kullanım Şartları
            </MuiLink>
          </Stack>
        </Stack>
      </Container>
    </Box>
  );
}

export default Footer;