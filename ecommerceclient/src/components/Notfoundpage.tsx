
import { Box, Button, Container, Typography, Paper, Stack, Grid } from '@mui/material';
import { useNavigate } from 'react-router';
import HomeRoundedIcon from '@mui/icons-material/HomeRounded';
import ShoppingBagOutlinedIcon from '@mui/icons-material/ShoppingBagOutlined';
import ArrowBackRoundedIcon from '@mui/icons-material/ArrowBackRounded';
import ExploreOutlinedIcon from '@mui/icons-material/ExploreOutlined';
import HelpOutlineRoundedIcon from '@mui/icons-material/HelpOutlineRounded';

function Notfoundpage() {
  const navigate = useNavigate();

  const helpfulLinks = [
    {
      title: 'Tüm Ürünler',
      description: 'En son eklenen ürünlerimizi ve kampanyaları keşfedin.',
      path: '/products',
      icon: <ExploreOutlinedIcon color="primary" sx={{ fontSize: 28 }} />,
    },
    {
      title: 'Sepetiniz',
      description: 'Eklediğiniz ürünleri inceleyin ve siparişinizi tamamlayın.',
      path: '/cart',
      icon: <ShoppingBagOutlinedIcon color="secondary" sx={{ fontSize: 28 }} />,
    },
    {
      title: 'İletişim & Destek',
      description: 'Yardıma mı ihtiyacınız var? Bizimle iletişime geçin.',
      path: '/contact',
      icon: <HelpOutlineRoundedIcon color="info" sx={{ fontSize: 28 }} />,
    },
  ];

  return (
    <Box
      sx={{
        minHeight: '85vh',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        py: { xs: 6, md: 10 },
        bgcolor: 'background.default',
      }}
    >
      <Container maxWidth="md" sx={{ textAlign: 'center' }}>
        {/* 404 Numarası - Düz Tema Rengi (Primary) */}
        <Typography
          variant="h1"
          sx={{
            fontSize: { xs: '6.5rem', sm: '9.5rem', md: '11rem' },
            fontWeight: 900,
            lineHeight: 1,
            letterSpacing: '-0.04em',
            color: 'primary.main',
            mb: 2,
            userSelect: 'none',
          }}
        >
          404
        </Typography>

        {/* Başlık ve Açıklama */}
        <Typography
          variant="h3"
          component="h1"
          sx={{
            fontWeight: 800,
            mb: 2,
            fontSize: { xs: '1.75rem', sm: '2.25rem', md: '2.5rem' },
            color: 'text.primary',
          }}
        >
          Sayfa Bulunamadı
        </Typography>

        <Typography
          variant="body1"
          sx={{
            color: 'text.secondary',
            maxWidth: 560,
            mx: 'auto',
            mb: 4.5,
            fontSize: { xs: '1rem', sm: '1.125rem' },
            lineHeight: 1.6,
          }}
        >
          Aradığınız sayfa silinmiş, adı değiştirilmiş veya geçici olarak kullanım dışı kalmış olabilir.
          Aşağıdaki butonlarla alışverişe kaldığınız yerden devam edebilirsiniz.
        </Typography>

        {/* Eylem Butonları - Düz Tema Renkleri */}
        <Stack
          direction={{ xs: 'column', sm: 'row' }}
          spacing={2}
          justifyContent="center"
          alignItems="center"
          sx={{ mb: 6 }}
        >
          <Button
            variant="contained"
            color="primary"
            size="large"
            startIcon={<HomeRoundedIcon />}
            onClick={() => navigate('/')}
            sx={{
              px: 3.5,
              py: 1.4,
              fontSize: '1rem',
              fontWeight: 600,
              borderRadius: '32px',
              transition: 'all 0.2s ease',
              '&:hover': {
                transform: 'translateY(-2px)',
              },
            }}
          >
            Ana Sayfaya Dön
          </Button>

          <Button
            variant="outlined"
            color="primary"
            size="large"
            startIcon={<ShoppingBagOutlinedIcon />}
            onClick={() => navigate('/products')}
            sx={{
              px: 3.5,
              py: 1.4,
              fontSize: '1rem',
              fontWeight: 600,
              borderRadius: '32px',
              transition: 'all 0.2s ease',
              '&:hover': {
                transform: 'translateY(-2px)',
              },
            }}
          >
            Ürünleri Keşfet
          </Button>

          <Button
            variant="text"
            size="large"
            startIcon={<ArrowBackRoundedIcon />}
            onClick={() => navigate(-1)}
            sx={{
              px: 2.5,
              py: 1.4,
              fontSize: '1rem',
              fontWeight: 500,
              color: 'text.secondary',
              borderRadius: '32px',
              '&:hover': {
                color: 'text.primary',
                bgcolor: 'action.hover',
              },
            }}
          >
            Önceki Sayfa
          </Button>
        </Stack>

        {/* Hızlı Öneri Kartları */}
        <Grid container spacing={2.5} sx={{ mt: 2 }}>
          {helpfulLinks.map((item, idx) => (
            <Grid key={idx} size={{ xs: 12, sm: 4 }}>
              <Paper
                elevation={0}
                onClick={() => navigate(item.path)}
                sx={{
                  p: 3,
                  height: '100%',
                  borderRadius: 2,
                  border: '1px solid',
                  borderColor: 'divider',
                  bgcolor: 'background.paper',
                  cursor: 'pointer',
                  textAlign: 'left',
                  display: 'flex',
                  flexDirection: 'column',
                  gap: 1.2,
                  transition: 'all 0.25s ease',
                  '&:hover': {
                    transform: 'translateY(-4px)',
                    borderColor: 'primary.main',
                    boxShadow: 2,
                  },
                }}
              >
                <Box sx={{ display: 'flex', alignItems: 'center', gap: 1.5 }}>
                  <Box
                    sx={{
                      p: 1,
                      borderRadius: 1.5,
                      bgcolor: 'grey.100',
                      display: 'flex',
                      alignItems: 'center',
                      justifyContent: 'center',
                    }}
                  >
                    {item.icon}
                  </Box>
                  <Typography variant="h6" sx={{ fontSize: '1.05rem', fontWeight: 700, color: 'text.primary' }}>
                    {item.title}
                  </Typography>
                </Box>
                <Typography variant="body2" color="text.secondary" sx={{ fontSize: '0.875rem' }}>
                  {item.description}
                </Typography>
              </Paper>
            </Grid>
          ))}
        </Grid>
      </Container>
    </Box>
  );
}

export default Notfoundpage;