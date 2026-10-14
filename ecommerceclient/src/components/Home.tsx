import React, { useEffect, useState } from 'react';
import {
  Container,
  Typography,
  Box,
  Button,
  IconButton,
  Grid,
  Paper,
  Rating,
  OutlinedInput,
  Avatar,
  Tooltip,
} from '@mui/material';
import { useNavigate } from 'react-router';
import { toast } from 'react-toastify';
import ShoppingBagOutlinedIcon from '@mui/icons-material/ShoppingBagOutlined';
import ShoppingCartOutlinedIcon from '@mui/icons-material/ShoppingCartOutlined';
import ArrowForwardIcon from '@mui/icons-material/ArrowForward';
import FavoriteBorderIcon from '@mui/icons-material/FavoriteBorder';
import FavoriteIcon from '@mui/icons-material/Favorite';
import CheckIcon from '@mui/icons-material/Check';
import LocalShippingOutlinedIcon from '@mui/icons-material/LocalShippingOutlined';
import VerifiedUserOutlinedIcon from '@mui/icons-material/VerifiedUserOutlined';
import AutorenewOutlinedIcon from '@mui/icons-material/AutorenewOutlined';
import SupportAgentOutlinedIcon from '@mui/icons-material/SupportAgentOutlined';
import SecurityOutlinedIcon from '@mui/icons-material/SecurityOutlined';
import CheckroomOutlinedIcon from '@mui/icons-material/CheckroomOutlined';
import HomeOutlinedIcon from '@mui/icons-material/HomeOutlined';
import DevicesOutlinedIcon from '@mui/icons-material/DevicesOutlined';
import SpaOutlinedIcon from '@mui/icons-material/SpaOutlined';
import FitnessCenterOutlinedIcon from '@mui/icons-material/FitnessCenterOutlined';
import WatchOutlinedIcon from '@mui/icons-material/WatchOutlined';
import AutoAwesomeOutlinedIcon from '@mui/icons-material/AutoAwesomeOutlined';
import MailOutlineRoundedIcon from '@mui/icons-material/MailOutlineRounded';
import VerifiedIcon from '@mui/icons-material/Verified';

import type { ProductListProps } from '../models/Props';
import { productService } from '../Services/productService';
import type { GetProducts } from '../models/ProductResponse';
import Spinner from './Spinner';

export default function Home(props: ProductListProps) {
  const navigate = useNavigate();
  const [products, setProducts] = useState<GetProducts[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [activeTab, setActiveTab] = useState<string>('all');
  const [favorites, setFavorites] = useState<Record<string, boolean>>({});
  const [newsletterEmail, setNewsletterEmail] = useState('');
  const [featuredIndex, setFeaturedIndex] = useState(0);

  // Countdown timer state (hours: 14, minutes: 48, seconds: 25)
  const [timeLeft, setTimeLeft] = useState({
    hours: 14,
    minutes: 48,
    seconds: 25,
  });

  // Countdown ticking effect
  useEffect(() => {
    const timer = setInterval(() => {
      setTimeLeft((prev) => {
        if (prev.seconds > 0) {
          return { ...prev, seconds: prev.seconds - 1 };
        } else if (prev.minutes > 0) {
          return { ...prev, minutes: prev.minutes - 1, seconds: 59 };
        } else if (prev.hours > 0) {
          return { hours: prev.hours - 1, minutes: 59, seconds: 59 };
        }
        return prev;
      });
    }, 1000);
    return () => clearInterval(timer);
  }, []);

  // Fetch real products from backend
  useEffect(() => {
    const fetchProducts = async () => {
      setLoading(true);
      try {
        const response = await productService.getPublicProducts(0, 50);
        if (response && response.Data && response.Data.length > 0) {
          setProducts(response.Data);
        }
      } catch (err) {
        console.error('Ürünler yüklenirken hata oluştu:', err);
      } finally {
        setLoading(false);
      }
    };
    fetchProducts();
  }, []);

  const toggleFavorite = (productId: string, e: React.MouseEvent) => {
    e.stopPropagation();
    setFavorites((prev) => {
      const newState = !prev[productId];
      if (newState) {
        toast.info('Ürün favorilere eklendi!');
      }
      return { ...prev, [productId]: newState };
    });
  };

  const handleNewsletterSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!newsletterEmail.trim() || !newsletterEmail.includes('@')) {
      toast.warning('Lütfen geçerli bir e-posta adresi girin.');
      return;
    }
    toast.success('Bülten kaydınız başarıyla alındı! 100 TL indirim kodunuz: HOSGELDIN100');
    setNewsletterEmail('');
  };

  if (loading) {
    return <Spinner message="Sayfa yükleniyor..." />;
  }

  // Active featured product for the Hero section
  const featuredProduct: GetProducts | undefined =
    products.length > 0 ? products[featuredIndex % products.length] : undefined;

  // Filtered products based on tabs
  const getFilteredProducts = () => {
    if (!products || products.length === 0) return [];
    switch (activeTab) {
      case 'new':
        return [...products].reverse();
      case 'deals':
        return [...products].sort((a, b) => a.price - b.price);
      case 'popular':
        return [...products].sort((a, b) => b.stock - a.stock);
      default:
        return products;
    }
  };

  const displayedProducts = getFilteredProducts();

  // Popular category items
  const popularCategories = [
    {
      id: 'fashion',
      title: 'Moda & Giyim',
      count: '450+ Ürün',
      icon: <CheckroomOutlinedIcon sx={{ fontSize: 28 }} />,
      query: 'Giyim',
    },
    {
      id: 'home-living',
      title: 'Ev & Yaşam',
      count: '320+ Ürün',
      icon: <HomeOutlinedIcon sx={{ fontSize: 28 }} />,
      query: 'Ev',
    },
    {
      id: 'electronics',
      title: 'Elektronik & Teknoloji',
      count: '210+ Ürün',
      icon: <DevicesOutlinedIcon sx={{ fontSize: 28 }} />,
      query: 'Elektronik',
    },
    {
      id: 'cosmetics',
      title: 'Kozmetik & Bakım',
      count: '180+ Ürün',
      icon: <SpaOutlinedIcon sx={{ fontSize: 28 }} />,
      query: 'Kozmetik',
    },
    {
      id: 'sports',
      title: 'Spor & Outdoor',
      count: '140+ Ürün',
      icon: <FitnessCenterOutlinedIcon sx={{ fontSize: 28 }} />,
      query: 'Spor',
    },
    {
      id: 'accessories',
      title: 'Aksesuar & Çanta',
      count: '290+ Ürün',
      icon: <WatchOutlinedIcon sx={{ fontSize: 28 }} />,
      query: 'Aksesuar',
    },
  ];

  // Testimonials
  const testimonials = [
    {
      id: 1,
      name: 'Emre Kaya',
      location: 'Onaylı Müşteri - İstanbul',
      avatar: 'EK',
      rating: 5,
      comment:
        'Siparişim ertesi gün kapımdaydı. Paketleme son derece özenli ve ürünler kusursuz geldi. Artık tüm online alışverişlerimde ilk tercihim burası.',
    },
    {
      id: 2,
      name: 'Selin Yılmaz',
      location: 'Onaylı Müşteri - İzmir',
      avatar: 'SY',
      rating: 5,
      comment:
        'Sipariş ettiğim ürünlerin kalitesi ve özeni birinci sınıf. Ayrıca müşteri hizmetleri sorularıma dakikalar içinde çok nazik ve ilgili yaklaştı.',
    },
    {
      id: 3,
      name: 'Burak Demir',
      location: 'Onaylı Müşteri - Ankara',
      avatar: 'BD',
      rating: 5,
      comment:
        'Fiyatlar piyasaya göre son derece avantajlı. Hem hızlı kargo süreci hem de 30 gün kolay iade politikası alışveriş yaparken büyük güven veriyor.',
    },
  ];

  // Value propositions
  const valueProps = [
    {
      icon: <LocalShippingOutlinedIcon sx={{ fontSize: 26, color: '#1E3A8A' }} />,
      title: 'Ücretsiz & Hızlı Kargo',
      desc: '150 TL ve üzeri siparişlerinizde aynı gün özenle paketlenip kargolanır.',
    },
    {
      icon: <AutorenewOutlinedIcon sx={{ fontSize: 26, color: '#1E3A8A' }} />,
      title: '30 Gün Kolay İade',
      desc: 'Sorunsuz iade süreci ve anında geri ödeme güvencesiyle risk almayın.',
    },
    {
      icon: <SupportAgentOutlinedIcon sx={{ fontSize: 26, color: '#1E3A8A' }} />,
      title: '7/24 Müşteri Desteği',
      desc: 'Alışveriş öncesi ve sonrasında canlı destek ekibimiz daima yanınızda.',
    },
    {
      icon: <SecurityOutlinedIcon sx={{ fontSize: 26, color: '#1E3A8A' }} />,
      title: '%100 Güvenli Ödeme',
      desc: '256-bit SSL şifreleme ve 3D Secure altyapısıyla güvenceli alışveriş yapın.',
    },
  ];

  return (
    <Box sx={{ bgcolor: 'background.default', pb: 8 }}>
      {/* ---------------- 1. HERO SECTION ---------------- */}
      <Container maxWidth="lg" sx={{ pt: { xs: 3, md: 5 }, pb: { xs: 4, md: 6 } }}>
        <Paper
          elevation={0}
          sx={{
            p: { xs: 3, sm: 4, md: 6 },
            borderRadius: 4,
            bgcolor: '#FFFFFF',
            border: '1px solid rgba(229, 231, 235, 0.8)',
            boxShadow: '0 1px 3px 0 rgba(0, 0, 0, 0.05), 0 1px 2px -1px rgba(0, 0, 0, 0.05)',
          }}
        >
          <Grid container spacing={{ xs: 4, md: 6 }} alignItems="center">
            {/* Sol Taraf: Metinler, Butonlar ve Güven Rozetleri */}
            <Grid size={{ xs: 12, md: 6 }}>
              <Box sx={{ maxWidth: 540 }}>
                {/* Üst Vurgu Rozeti */}
                <Box
                  sx={{
                    display: 'inline-flex',
                    alignItems: 'center',
                    gap: 0.75,
                    px: 1.75,
                    py: 0.6,
                    borderRadius: 9999,
                    bgcolor: '#EFF6FF',
                    border: '1px solid #BFDBFE',
                    color: '#1E3A8A',
                    mb: 2.5,
                  }}
                >
                  <AutoAwesomeOutlinedIcon sx={{ fontSize: 18, color: '#F97316' }} />
                  <Typography
                    variant="caption"
                    sx={{
                      fontFamily: '"Outfit", sans-serif',
                      fontWeight: 700,
                      fontSize: '0.8rem',
                      letterSpacing: '0.3px',
                    }}
                  >
                    Seçkin Koleksiyonlar & Yeni Sezon
                  </Typography>
                </Box>

                {/* Ana Başlık */}
                <Typography
                  variant="h1"
                  component="h1"
                  sx={{
                    fontFamily: '"Outfit", sans-serif',
                    fontSize: { xs: '2.2rem', sm: '2.8rem', md: '3.4rem' },
                    fontWeight: 800,
                    lineHeight: 1.15,
                    color: 'primary.main',
                    letterSpacing: '-1px',
                    mb: 2,
                  }}
                >
                  Yaşamınıza Değer Katan{' '}
                  <Box component="span" sx={{ color: 'secondary.main' }}>
                    Özel
                  </Box>
                  <br />
                  <Box component="span" sx={{ color: 'secondary.main' }}>
                    Koleksiyonlar
                  </Box>
                </Typography>

                {/* Alt Metin */}
                <Typography
                  variant="body1"
                  sx={{
                    fontFamily: '"Inter", sans-serif',
                    color: 'text.secondary',
                    fontSize: { xs: '0.95rem', md: '1.05rem' },
                    lineHeight: 1.65,
                    mb: 4,
                  }}
                >
                  Modadan ev yaşamına, teknolojiden aksesuara en trend ürünler, avantajlı fiyatlar ve güvenilir
                  teslimat ile kapınızda. Keyifli ve güvenli alışverişi keşfedin.
                </Typography>

                {/* Aksiyon Butonları */}
                <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 2, mb: 4 }}>
                  <Button
                    variant="contained"
                    size="large"
                    startIcon={<ShoppingBagOutlinedIcon />}
                    onClick={() => navigate('/products')}
                    sx={{
                      bgcolor: 'secondary.main',
                      color: '#FFFFFF',
                      px: 3.5,
                      py: 1.4,
                      fontSize: '1rem',
                      fontWeight: 700,
                      borderRadius: '10px',
                      boxShadow: '0 4px 12px rgba(249, 115, 22, 0.25)',
                      '&:hover': {
                        bgcolor: 'secondary.dark',
                        boxShadow: '0 6px 18px rgba(249, 115, 22, 0.35)',
                      },
                    }}
                  >
                    Alışverişe Başla
                  </Button>

                  <Button
                    variant="outlined"
                    size="large"
                    endIcon={<ArrowForwardIcon />}
                    onClick={() => navigate('/products')}
                    sx={{
                      bgcolor: '#FFFFFF',
                      borderColor: '#D1D5DB',
                      color: 'text.primary',
                      px: 3,
                      py: 1.4,
                      fontSize: '1rem',
                      fontWeight: 600,
                      borderRadius: '10px',
                      '&:hover': {
                        borderColor: 'primary.main',
                        color: 'primary.main',
                        bgcolor: 'rgba(30, 58, 138, 0.04)',
                      },
                    }}
                  >
                    Öne Çıkan Ürünleri Gör
                  </Button>
                </Box>

                {/* Rozetler / Avantaj Maddeleri */}
                <Box
                  sx={{
                    display: 'flex',
                    flexWrap: 'wrap',
                    alignItems: 'center',
                    gap: { xs: 2, sm: 3 },
                    pt: 2,
                    borderTop: '1px solid #F3F4F6',
                  }}
                >
                  <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                    <LocalShippingOutlinedIcon sx={{ fontSize: 20, color: 'primary.light' }} />
                    <Typography variant="caption" sx={{ fontWeight: 600, color: 'text.secondary' }}>
                      Aynı Gün Kargo
                    </Typography>
                  </Box>
                  <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                    <VerifiedUserOutlinedIcon sx={{ fontSize: 20, color: 'primary.light' }} />
                    <Typography variant="caption" sx={{ fontWeight: 600, color: 'text.secondary' }}>
                      Orijinal Ürün
                    </Typography>
                  </Box>
                  <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                    <AutorenewOutlinedIcon sx={{ fontSize: 20, color: 'primary.light' }} />
                    <Typography variant="caption" sx={{ fontWeight: 600, color: 'text.secondary' }}>
                      30 Gün Kolay İade
                    </Typography>
                  </Box>
                </Box>
              </Box>
            </Grid>

            {/* Sağ Taraf: Hero Ürün Vitrin Kartı (Mevcut Ürünlerden Çekilir) */}
            <Grid size={{ xs: 12, md: 6 }}>
              {featuredProduct ? (
                <Box
                  sx={{
                    position: 'relative',
                    width: '100%',
                    maxWidth: { xs: '100%', md: 460 },
                    mx: 'auto',
                    p: { xs: 1.5, sm: 2 },
                    borderRadius: 4,
                    bgcolor: '#FFFFFF',
                    border: '1px solid rgba(229, 231, 235, 0.8)',
                    boxShadow: '0 10px 25px -5px rgba(30, 58, 138, 0.08)',
                    cursor: 'pointer',
                    transition: 'all 0.35s ease',
                    '&:hover': {
                      transform: 'translateY(-4px)',
                      boxShadow: '0 20px 35px -8px rgba(30, 58, 138, 0.15)',
                      '& .hero-product-img': {
                        transform: 'scale(1.04)',
                      },
                    },
                  }}
                  onClick={() => navigate(`/product/${featuredProduct.product_id}`)}
                >
                  {/* Ürün Görseli - Yuvarlatılmış (Border Radius) Çerçeve */}
                  <Box
                    sx={{
                      position: 'relative',
                      width: '100%',
                      aspectRatio: '1 / 1',
                      maxHeight: { xs: 320, sm: 380, md: 430 },
                      borderRadius: 3,
                      overflow: 'hidden',
                      bgcolor: '#F3F4F6',
                      display: 'flex',
                      alignItems: 'center',
                      justifyContent: 'center',
                    }}
                  >
                    {/* Sol Üst İndirim Rozeti */}
                    <Box
                      sx={{
                        position: 'absolute',
                        top: 12,
                        left: 12,
                        zIndex: 3,
                        bgcolor: 'secondary.main',
                        color: 'white',
                        px: 1.5,
                        py: 0.5,
                        borderRadius: 9999,
                        fontSize: '0.75rem',
                        fontWeight: 800,
                        letterSpacing: '0.5px',
                        boxShadow: '0 4px 10px rgba(249, 115, 22, 0.3)',
                      }}
                    >
                      %25 İNDİRİM
                    </Box>

                    {/* Sağ Üst Favori Butonu */}
                    <Tooltip title="Favorilere Ekle">
                      <IconButton
                        onClick={(e) => toggleFavorite(featuredProduct.product_id, e)}
                        sx={{
                          position: 'absolute',
                          top: 10,
                          right: 10,
                          zIndex: 3,
                          bgcolor: 'rgba(255, 255, 255, 0.9)',
                          backdropFilter: 'blur(6px)',
                          boxShadow: '0 2px 8px rgba(0, 0, 0, 0.1)',
                          '&:hover': { bgcolor: '#FFFFFF' },
                        }}
                      >
                        {favorites[featuredProduct.product_id] ? (
                          <FavoriteIcon sx={{ color: '#EF4444', fontSize: 20 }} />
                        ) : (
                          <FavoriteBorderIcon sx={{ color: 'text.secondary', fontSize: 20 }} />
                        )}
                      </IconButton>
                    </Tooltip>

                    {/* Fotoğraf (Köşeleri Yuvarlatılmış ve Çerçeveye Uyumlu) */}
                    <Box
                      component="img"
                      className="hero-product-img"
                      src={featuredProduct.thumbnail_image}
                      alt={featuredProduct.product_name}
                      sx={{
                        width: '100%',
                        height: '100%',
                        objectFit: 'cover',
                        borderRadius: 3,
                        transition: 'transform 0.4s cubic-bezier(0.2, 0, 0.2, 1)',
                      }}
                    />
                  </Box>

                  {/* Alt Kısım: Ürün Detay Kartı */}
                  <Box
                    sx={{
                      pt: 2,
                      pb: 0.5,
                      px: 0.5,
                    }}
                  >
                    <Typography
                      variant="caption"
                      sx={{
                        fontFamily: '"Outfit", sans-serif',
                        fontWeight: 700,
                        color: 'primary.light',
                        letterSpacing: '1px',
                        textTransform: 'uppercase',
                        fontSize: '0.72rem',
                        display: 'block',
                        mb: 0.5,
                      }}
                    >
                      GÜNÜN FIRSATI
                    </Typography>

                    <Box
                      sx={{
                        display: 'flex',
                        alignItems: 'center',
                        justifyContent: 'space-between',
                        mb: 0.75,
                      }}
                    >
                      <Typography
                        variant="h6"
                        sx={{
                          fontFamily: '"Outfit", sans-serif',
                          fontWeight: 700,
                          color: 'primary.main',
                          fontSize: '1.15rem',
                          lineHeight: 1.2,
                        }}
                      >
                        {featuredProduct.product_name}
                      </Typography>

                      <Typography
                        variant="body2"
                        sx={{
                          color: 'text.disabled',
                          textDecoration: 'line-through',
                          fontSize: '0.85rem',
                          fontWeight: 500,
                        }}
                      >
                        {(featuredProduct.price * 1.33).toFixed(0)} ₺
                      </Typography>
                    </Box>

                    <Box
                      sx={{
                        display: 'flex',
                        alignItems: 'center',
                        justifyContent: 'space-between',
                      }}
                    >
                      <Box sx={{ display: 'flex', alignItems: 'center', gap: 0.75 }}>
                        <Rating value={5} readOnly size="small" sx={{ fontSize: '1rem', color: '#F59E0B' }} />
                        <Typography variant="caption" sx={{ color: 'text.secondary', fontWeight: 600 }}>
                          4.9 (120+ Değerlendirme)
                        </Typography>
                      </Box>

                      <Typography
                        variant="h6"
                        sx={{
                          fontFamily: '"Outfit", sans-serif',
                          fontWeight: 800,
                          color: 'primary.main',
                          fontSize: '1.25rem',
                        }}
                      >
                        {featuredProduct.price.toLocaleString('tr-TR')} ₺
                      </Typography>
                    </Box>

                    {/* Eğer birden fazla ürün varsa vitrini değiştirebilme butonları */}
                    {products.length > 1 && (
                      <Box sx={{ display: 'flex', gap: 1, mt: 1.5, justifyContent: 'center' }}>
                        {products.slice(0, 4).map((p, idx) => (
                          <Box
                            key={p.product_id}
                            onClick={(e) => {
                              e.stopPropagation();
                              setFeaturedIndex(idx);
                            }}
                            sx={{
                              width: idx === featuredIndex ? 24 : 8,
                              height: 8,
                              borderRadius: 4,
                              bgcolor: idx === featuredIndex ? 'primary.main' : '#D1D5DB',
                              cursor: 'pointer',
                              transition: 'all 0.25s ease',
                            }}
                          />
                        ))}
                      </Box>
                    )}
                  </Box>
                </Box>
              ) : null}
            </Grid>
          </Grid>
        </Paper>
      </Container>

      {/* ---------------- 2. POPÜLER KATEGORİLER ---------------- */}
      <Container maxWidth="lg" sx={{ py: { xs: 4, md: 5 } }}>
        <Box
          sx={{
            display: 'flex',
            flexDirection: { xs: 'column', sm: 'row' },
            alignItems: { xs: 'flex-start', sm: 'center' },
            justifyContent: 'space-between',
            mb: 3.5,
            gap: 1,
          }}
        >
          <Box>
            <Typography
              variant="h4"
              component="h2"
              sx={{
                fontFamily: '"Outfit", sans-serif',
                fontWeight: 700,
                color: 'primary.main',
                fontSize: { xs: '1.5rem', md: '1.85rem' },
                letterSpacing: '-0.5px',
              }}
            >
              Popüler Kategoriler
            </Typography>
            <Typography
              variant="body2"
              sx={{
                fontFamily: '"Inter", sans-serif',
                color: 'text.secondary',
                mt: 0.5,
              }}
            >
              İhtiyacınıza ve tarzınıza en uygun kaliteli ürünleri keşfedin
            </Typography>
          </Box>

          <Button
            variant="text"
            endIcon={<ArrowForwardIcon />}
            onClick={() => navigate('/products')}
            sx={{
              fontFamily: '"Outfit", sans-serif',
              fontWeight: 700,
              color: 'primary.main',
              fontSize: '0.95rem',
              p: 0.5,
              '&:hover': {
                bgcolor: 'transparent',
                color: 'primary.light',
                transform: 'translateX(3px)',
              },
              transition: 'all 0.2s ease',
            }}
          >
            Tümünü Gör
          </Button>
        </Box>

        <Grid container spacing={{ xs: 2, sm: 2.5, md: 3 }}>
          {popularCategories.map((cat) => (
            <Grid size={{ xs: 6, sm: 4, md: 2 }} key={cat.id}>
              <Paper
                elevation={0}
                onClick={() => navigate(`/products?searchQuery=${encodeURIComponent(cat.query)}`)}
                sx={{
                  p: { xs: 2, md: 2.5 },
                  borderRadius: 3,
                  bgcolor: '#FFFFFF',
                  border: '1px solid rgba(229, 231, 235, 0.8)',
                  boxShadow: '0 1px 3px 0 rgba(0, 0, 0, 0.05), 0 1px 2px -1px rgba(0, 0, 0, 0.05)',
                  textAlign: 'center',
                  cursor: 'pointer',
                  transition: 'all 0.25s cubic-bezier(0.16, 1, 0.3, 1)',
                  display: 'flex',
                  flexDirection: 'column',
                  alignItems: 'center',
                  justifyContent: 'center',
                  minHeight: 140,
                  '&:hover': {
                    transform: 'translateY(-5px)',
                    boxShadow: '0 10px 15px -3px rgba(30, 58, 138, 0.08), 0 4px 6px -4px rgba(30, 58, 138, 0.04)',
                    borderColor: '#BFDBFE',
                    '& .category-icon-box': {
                      bgcolor: 'primary.main',
                      color: '#FFFFFF',
                    },
                  },
                }}
              >
                <Box
                  className="category-icon-box"
                  sx={{
                    width: 52,
                    height: 52,
                    borderRadius: '50%',
                    bgcolor: '#EFF6FF',
                    color: 'primary.main',
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'center',
                    mb: 1.5,
                    transition: 'all 0.25s ease',
                  }}
                >
                  {cat.icon}
                </Box>
                <Typography
                  variant="subtitle2"
                  sx={{
                    fontFamily: '"Outfit", sans-serif',
                    fontWeight: 700,
                    color: 'primary.main',
                    fontSize: '0.9rem',
                    lineHeight: 1.2,
                    mb: 0.5,
                  }}
                >
                  {cat.title}
                </Typography>
                <Typography
                  variant="caption"
                  sx={{
                    fontFamily: '"Inter", sans-serif',
                    color: 'text.secondary',
                    fontSize: '0.75rem',
                  }}
                >
                  {cat.count}
                </Typography>
              </Paper>
            </Grid>
          ))}
        </Grid>
      </Container>

      {/* ---------------- 3. ÇOK SATANLAR & TREND ÜRÜNLER ---------------- */}
      <Container maxWidth="lg" sx={{ py: { xs: 4, md: 6 } }}>
        <Box
          sx={{
            display: 'flex',
            flexDirection: { xs: 'column', md: 'row' },
            alignItems: { xs: 'flex-start', md: 'center' },
            justifyContent: 'space-between',
            mb: 4,
            gap: 2,
          }}
        >
          <Box>
            <Typography
              variant="h4"
              component="h2"
              sx={{
                fontFamily: '"Outfit", sans-serif',
                fontWeight: 700,
                color: 'primary.main',
                fontSize: { xs: '1.5rem', md: '1.85rem' },
                letterSpacing: '-0.5px',
              }}
            >
              Çok Satanlar & Trend Ürünler
            </Typography>
            <Typography
              variant="body2"
              sx={{
                fontFamily: '"Inter", sans-serif',
                color: 'text.secondary',
                mt: 0.5,
              }}
            >
              Kullanıcılarımızın en çok ilgi gösterdiği ve beğendiği popüler ürünler.
            </Typography>
          </Box>

          {/* Filtre Butonları (Pill Şeklinde) */}
          <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 1 }}>
            {[
              { id: 'all', label: 'Tümü' },
              { id: 'new', label: 'Yeni Çıkanlar' },
              { id: 'deals', label: 'Fırsat Ürünleri' },
              { id: 'popular', label: 'En Çok Değerlendirilenler' },
            ].map((tab) => (
              <Button
                key={tab.id}
                size="small"
                onClick={() => setActiveTab(tab.id)}
                sx={{
                  borderRadius: 9999,
                  px: 2,
                  py: 0.6,
                  fontSize: '0.825rem',
                  fontWeight: 600,
                  bgcolor: activeTab === tab.id ? 'primary.main' : '#FFFFFF',
                  color: activeTab === tab.id ? '#FFFFFF' : 'text.secondary',
                  border: '1px solid',
                  borderColor: activeTab === tab.id ? 'primary.main' : '#E5E7EB',
                  '&:hover': {
                    bgcolor: activeTab === tab.id ? 'primary.dark' : '#F3F4F6',
                  },
                }}
              >
                {tab.label}
              </Button>
            ))}
          </Box>
        </Box>

        {/* 4 Kolonlu Ürün Kartları Grid */}
        <Grid container spacing={{ xs: 2.5, sm: 3, md: 3.5 }}>
          {displayedProducts.map((product, index) => {
            const inCart = props.cartItems
              ? props.cartItems.some((item) => item.product_id === product.product_id)
              : false;
            const isFav = !!favorites[product.product_id];
            // Dinamik rozetler
            const badgeText = index === 0 ? '%15 İNDİRİM' : index === 2 ? '%20 İNDİRİM' : null;
            const tagText = index % 2 === 0 ? 'Ücretsiz Kargo' : 'Hızlı Gönderi';

            return (
              <Grid size={{ xs: 12, sm: 6, md: 3 }} key={product.product_id}>
                <Paper
                  elevation={0}
                  onClick={() => navigate(`/product/${product.product_id}`)}
                  sx={{
                    height: '100%',
                    display: 'flex',
                    flexDirection: 'column',
                    borderRadius: 3,
                    bgcolor: '#FFFFFF',
                    border: '1px solid rgba(229, 231, 235, 0.8)',
                    boxShadow: '0 1px 3px 0 rgba(0, 0, 0, 0.05), 0 1px 2px -1px rgba(0, 0, 0, 0.05)',
                    overflow: 'hidden',
                    cursor: 'pointer',
                    transition: 'all 0.3s cubic-bezier(0.16, 1, 0.3, 1)',
                    '&:hover': {
                      transform: 'translateY(-6px)',
                      boxShadow: '0 12px 24px -4px rgba(30, 58, 138, 0.12)',
                      borderColor: '#BFDBFE',
                      '& .product-card-img': {
                        transform: 'scale(1.06)',
                      },
                    },
                  }}
                >
                  {/* Görsel Alanı */}
                  <Box
                    sx={{
                      position: 'relative',
                      height: 220,
                      bgcolor: '#FFFFFF',
                      display: 'flex',
                      alignItems: 'center',
                      justifyContent: 'center',
                      p: 2.5,
                      borderBottom: '1px solid #F3F4F6',
                    }}
                  >
                    {/* Sol Üst İndirim Rozeti */}
                    {badgeText && (
                      <Box
                        sx={{
                          position: 'absolute',
                          top: 12,
                          left: 12,
                          bgcolor: 'secondary.main',
                          color: 'white',
                          px: 1,
                          py: 0.3,
                          borderRadius: 1.5,
                          fontSize: '0.72rem',
                          fontWeight: 800,
                          letterSpacing: '0.3px',
                          boxShadow: '0 2px 6px rgba(249, 115, 22, 0.25)',
                          zIndex: 2,
                        }}
                      >
                        {badgeText}
                      </Box>
                    )}

                    {/* Sağ Üst Favori Butonu */}
                    <Tooltip title={isFav ? 'Favorilerden Çıkar' : 'Favorilere Ekle'}>
                      <IconButton
                        onClick={(e) => toggleFavorite(product.product_id, e)}
                        sx={{
                          position: 'absolute',
                          top: 10,
                          right: 10,
                          zIndex: 2,
                          bgcolor: 'rgba(255, 255, 255, 0.85)',
                          backdropFilter: 'blur(4px)',
                          boxShadow: '0 2px 6px rgba(0, 0, 0, 0.08)',
                          '&:hover': { bgcolor: '#FFFFFF' },
                        }}
                      >
                        {isFav ? (
                          <FavoriteIcon sx={{ color: '#EF4444', fontSize: 18 }} />
                        ) : (
                          <FavoriteBorderIcon sx={{ color: 'text.secondary', fontSize: 18 }} />
                        )}
                      </IconButton>
                    </Tooltip>

                    {/* Ürün Fotoğrafı */}
                    <Box
                      component="img"
                      className="product-card-img"
                      src={product.thumbnail_image}
                      alt={product.product_name}
                      sx={{
                        maxWidth: '100%',
                        maxHeight: '100%',
                        objectFit: 'contain',
                        transition: 'transform 0.4s ease',
                      }}
                    />
                  </Box>

                  {/* Ürün Detayları */}
                  <Box
                    sx={{
                      p: 2.2,
                      display: 'flex',
                      flexDirection: 'column',
                      flexGrow: 1,
                    }}
                  >
                    {/* Kategori Alt Başlığı */}
                    <Typography
                      variant="caption"
                      sx={{
                        fontFamily: '"Inter", sans-serif',
                        fontWeight: 700,
                        color: '#9CA3AF',
                        textTransform: 'uppercase',
                        letterSpacing: '0.8px',
                        fontSize: '0.7rem',
                        mb: 0.5,
                      }}
                    >
                      ÖNE ÇIKAN KOLEKSİYON
                    </Typography>

                    {/* Ürün Başlığı */}
                    <Typography
                      variant="subtitle1"
                      sx={{
                        fontFamily: '"Outfit", sans-serif',
                        fontWeight: 700,
                        color: 'primary.main',
                        fontSize: '0.975rem',
                        lineHeight: 1.3,
                        mb: 1,
                        display: '-webkit-box',
                        WebkitLineClamp: 2,
                        WebkitBoxOrient: 'vertical',
                        overflow: 'hidden',
                        height: 40,
                      }}
                    >
                      {product.product_name}
                    </Typography>

                    {/* Yıldız Değerlendirme */}
                    <Box sx={{ display: 'flex', alignItems: 'center', gap: 0.5, mb: 1.5 }}>
                      <Rating
                        value={4.8}
                        readOnly
                        precision={0.1}
                        size="small"
                        sx={{ fontSize: '0.95rem', color: '#F59E0B' }}
                      />
                      <Typography variant="caption" sx={{ color: 'text.disabled', fontSize: '0.75rem', fontWeight: 600 }}>
                        (342)
                      </Typography>
                    </Box>

                    {/* Fiyat ve Etiket */}
                    <Box
                      sx={{
                        display: 'flex',
                        alignItems: 'baseline',
                        justifyContent: 'space-between',
                        mb: 2,
                        mt: 'auto',
                      }}
                    >
                      <Box>
                        <Typography
                          variant="caption"
                          sx={{
                            color: 'text.disabled',
                            textDecoration: 'line-through',
                            fontSize: '0.75rem',
                            display: 'block',
                            lineHeight: 1,
                            mb: 0.3,
                          }}
                        >
                          {(product.price * 1.2).toFixed(0)} ₺
                        </Typography>
                        <Typography
                          variant="h6"
                          sx={{
                            fontFamily: '"Outfit", sans-serif',
                            fontWeight: 800,
                            color: 'primary.main',
                            fontSize: '1.2rem',
                            lineHeight: 1,
                          }}
                        >
                          {product.price.toLocaleString('tr-TR')} ₺
                        </Typography>
                      </Box>

                      <Box
                        sx={{
                          px: 1,
                          py: 0.3,
                          borderRadius: 1,
                          bgcolor: '#EFF6FF',
                          border: '1px solid #DBEAFE',
                          color: 'primary.light',
                          fontSize: '0.7rem',
                          fontWeight: 600,
                        }}
                      >
                        {tagText}
                      </Box>
                    </Box>

                    {/* Sepete Ekle Butonu */}
                    <Button
                      fullWidth
                      variant="contained"
                      startIcon={inCart ? <CheckIcon /> : <ShoppingCartOutlinedIcon />}
                      disabled={inCart}
                      onClick={(e) => {
                        e.stopPropagation();
                        props.addtocart({
                          product_id: product.product_id,
                          product_name: product.product_name,
                          price: product.price,
                          stock: product.stock,
                          thumbnail_image: product.thumbnail_image,
                          quantity: 1,
                        });
                      }}
                      sx={{
                        borderRadius: '10px',
                        py: 1,
                        fontWeight: 700,
                        fontSize: '0.875rem',
                        bgcolor: inCart ? 'success.main' : 'primary.main',
                        color: '#FFFFFF',
                        '&:hover': {
                          bgcolor: inCart ? 'success.dark' : 'primary.dark',
                        },
                        '&.Mui-disabled': {
                          bgcolor: 'success.main',
                          color: '#FFFFFF',
                          opacity: 0.9,
                        },
                      }}
                    >
                      {inCart ? 'Sepete Eklendi' : 'Sepete Ekle'}
                    </Button>
                  </Box>
                </Paper>
              </Grid>
            );
          })}
        </Grid>
      </Container>

      {/* ---------------- 4. SINIRLI SÜRELİ HAFTA SONU FIRSATI BANNER ---------------- */}
      <Container maxWidth="lg" sx={{ py: { xs: 4, md: 6 } }}>
        <Paper
          elevation={0}
          sx={{
            borderRadius: 4,
            overflow: 'hidden',
            background: 'linear-gradient(135deg, #0F2042 0%, #172554 45%, #1E3A8A 100%)',
            color: '#FFFFFF',
            p: { xs: 3.5, sm: 5, md: 6 },
            position: 'relative',
            boxShadow: '0 12px 30px -5px rgba(23, 37, 84, 0.35)',
          }}
        >
          <Grid container spacing={4} alignItems="center" justifyContent="space-between">
            {/* Sol Kısım: Kampanya Bilgisi ve Canlı Geri Sayım */}
            <Grid size={{ xs: 12, md: 8 }}>
              {/* Kampanya Rozeti */}
              <Box
                sx={{
                  display: 'inline-flex',
                  alignItems: 'center',
                  gap: 1,
                  px: 1.75,
                  py: 0.6,
                  borderRadius: 9999,
                  bgcolor: 'rgba(255, 255, 255, 0.1)',
                  backdropFilter: 'blur(8px)',
                  border: '1px solid rgba(255, 255, 255, 0.2)',
                  mb: 2,
                }}
              >
                <Typography
                  variant="caption"
                  sx={{
                    fontFamily: '"Outfit", sans-serif',
                    fontWeight: 700,
                    fontSize: '0.78rem',
                    letterSpacing: '1px',
                    textTransform: 'uppercase',
                    color: '#93C5FD',
                  }}
                >
                  ⏱️ Sınırlı Süreli Hafta Sonu Fırsatı
                </Typography>
              </Box>

              {/* Başlık */}
              <Typography
                variant="h3"
                component="h2"
                sx={{
                  fontFamily: '"Outfit", sans-serif',
                  fontWeight: 800,
                  fontSize: { xs: '1.75rem', sm: '2.25rem', md: '2.6rem' },
                  lineHeight: 1.2,
                  letterSpacing: '-0.5px',
                  color: '#FFFFFF',
                  mb: 1.5,
                }}
              >
                Seçili Tüm Kategorilerde{' '}
                <Box component="span" sx={{ color: 'secondary.main' }}>
                  %30'a Varan
                </Box>{' '}
                İndirimler
              </Typography>

              {/* Açıklama */}
              <Typography
                variant="body1"
                sx={{
                  color: 'rgba(255, 255, 255, 0.85)',
                  fontSize: { xs: '0.95rem', md: '1.05rem' },
                  maxWidth: 600,
                  lineHeight: 1.6,
                  mb: 3.5,
                }}
              >
                Sezonun en gözde koleksiyonlarından günlük ihtiyaçlarınıza kadar binlerce üründe
                kaçırılmayacak fırsatlar devam ediyor.
              </Typography>

              {/* Geri Sayım Sayacı */}
              <Box sx={{ display: 'flex', alignItems: 'center', gap: { xs: 1.5, sm: 2 } }}>
                {/* Saat Kutusu */}
                <Box
                  sx={{
                    minWidth: { xs: 60, sm: 70 },
                    p: { xs: 1, sm: 1.5 },
                    borderRadius: 2.5,
                    bgcolor: 'rgba(11, 25, 56, 0.75)',
                    border: '1px solid rgba(255, 255, 255, 0.15)',
                    textAlign: 'center',
                  }}
                >
                  <Typography
                    variant="h5"
                    sx={{
                      fontFamily: '"Outfit", sans-serif',
                      fontWeight: 800,
                      color: '#FFFFFF',
                      fontSize: { xs: '1.4rem', sm: '1.75rem' },
                      lineHeight: 1,
                    }}
                  >
                    {String(timeLeft.hours).padStart(2, '0')}
                  </Typography>
                  <Typography
                    variant="caption"
                    sx={{
                      fontSize: '0.65rem',
                      fontWeight: 700,
                      color: '#93C5FD',
                      letterSpacing: '0.5px',
                      mt: 0.5,
                      display: 'block',
                    }}
                  >
                    SAAT
                  </Typography>
                </Box>

                <Typography variant="h5" sx={{ fontWeight: 800, color: 'rgba(255, 255, 255, 0.5)' }}>
                  :
                </Typography>

                {/* Dakika Kutusu */}
                <Box
                  sx={{
                    minWidth: { xs: 60, sm: 70 },
                    p: { xs: 1, sm: 1.5 },
                    borderRadius: 2.5,
                    bgcolor: 'rgba(11, 25, 56, 0.75)',
                    border: '1px solid rgba(255, 255, 255, 0.15)',
                    textAlign: 'center',
                  }}
                >
                  <Typography
                    variant="h5"
                    sx={{
                      fontFamily: '"Outfit", sans-serif',
                      fontWeight: 800,
                      color: '#FFFFFF',
                      fontSize: { xs: '1.4rem', sm: '1.75rem' },
                      lineHeight: 1,
                    }}
                  >
                    {String(timeLeft.minutes).padStart(2, '0')}
                  </Typography>
                  <Typography
                    variant="caption"
                    sx={{
                      fontSize: '0.65rem',
                      fontWeight: 700,
                      color: '#93C5FD',
                      letterSpacing: '0.5px',
                      mt: 0.5,
                      display: 'block',
                    }}
                  >
                    DAKİKA
                  </Typography>
                </Box>

                <Typography variant="h5" sx={{ fontWeight: 800, color: 'rgba(255, 255, 255, 0.5)' }}>
                  :
                </Typography>

                {/* Saniye Kutusu */}
                <Box
                  sx={{
                    minWidth: { xs: 60, sm: 70 },
                    p: { xs: 1, sm: 1.5 },
                    borderRadius: 2.5,
                    bgcolor: 'rgba(11, 25, 56, 0.75)',
                    border: '1px solid rgba(255, 255, 255, 0.15)',
                    textAlign: 'center',
                  }}
                >
                  <Typography
                    variant="h5"
                    sx={{
                      fontFamily: '"Outfit", sans-serif',
                      fontWeight: 800,
                      color: '#FFFFFF',
                      fontSize: { xs: '1.4rem', sm: '1.75rem' },
                      lineHeight: 1,
                    }}
                  >
                    {String(timeLeft.seconds).padStart(2, '0')}
                  </Typography>
                  <Typography
                    variant="caption"
                    sx={{
                      fontSize: '0.65rem',
                      fontWeight: 700,
                      color: '#93C5FD',
                      letterSpacing: '0.5px',
                      mt: 0.5,
                      display: 'block',
                    }}
                  >
                    SANİYE
                  </Typography>
                </Box>
              </Box>
            </Grid>

            {/* Sağ Kısım: CTA Butonu */}
            <Grid size={{ xs: 12, md: 4 }} sx={{ textAlign: { xs: 'left', md: 'right' } }}>
              <Button
                variant="contained"
                size="large"
                endIcon={<ArrowForwardIcon />}
                onClick={() => navigate('/products')}
                sx={{
                  bgcolor: 'secondary.main',
                  color: '#FFFFFF',
                  px: 4.5,
                  py: 1.6,
                  fontSize: '1.05rem',
                  fontWeight: 700,
                  borderRadius: '10px',
                  boxShadow: '0 6px 20px rgba(249, 115, 22, 0.4)',
                  '&:hover': {
                    bgcolor: 'secondary.dark',
                    boxShadow: '0 8px 25px rgba(249, 115, 22, 0.5)',
                  },
                }}
              >
                Fırsatları Yakala
              </Button>
              <Typography
                variant="caption"
                sx={{
                  display: 'block',
                  color: 'rgba(255, 255, 255, 0.6)',
                  mt: 1.5,
                  fontSize: '0.75rem',
                }}
              >
                *Kampanya stoklarla sınırlıdır ve sepette otomatik uygulanır.
              </Typography>
            </Grid>
          </Grid>
        </Paper>
      </Container>

      {/* ---------------- 5. DEĞER / GÜVEN ROZETLERİ ---------------- */}
      <Container maxWidth="lg" sx={{ py: { xs: 3, md: 5 } }}>
        <Grid container spacing={{ xs: 2, sm: 2.5, md: 3 }}>
          {valueProps.map((item, idx) => (
            <Grid size={{ xs: 12, sm: 6, md: 3 }} key={idx}>
              <Paper
                elevation={0}
                sx={{
                  p: 3,
                  height: '100%',
                  borderRadius: 3,
                  bgcolor: '#FFFFFF',
                  border: '1px solid rgba(229, 231, 235, 0.8)',
                  boxShadow: '0 1px 3px 0 rgba(0, 0, 0, 0.05), 0 1px 2px -1px rgba(0, 0, 0, 0.05)',
                  transition: 'all 0.25s ease',
                  '&:hover': {
                    boxShadow: '0 6px 15px rgba(30, 58, 138, 0.06)',
                    borderColor: '#BFDBFE',
                  },
                }}
              >
                <Box
                  sx={{
                    width: 46,
                    height: 46,
                    borderRadius: 2,
                    bgcolor: '#EFF6FF',
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'center',
                    mb: 1.5,
                  }}
                >
                  {item.icon}
                </Box>
                <Typography
                  variant="subtitle1"
                  sx={{
                    fontFamily: '"Outfit", sans-serif',
                    fontWeight: 700,
                    color: 'primary.main',
                    fontSize: '1rem',
                    lineHeight: 1.3,
                    mb: 0.5,
                  }}
                >
                  {item.title}
                </Typography>
                <Typography
                  variant="body2"
                  sx={{
                    fontFamily: '"Inter", sans-serif',
                    color: 'text.secondary',
                    fontSize: '0.825rem',
                    lineHeight: 1.5,
                  }}
                >
                  {item.desc}
                </Typography>
              </Paper>
            </Grid>
          ))}
        </Grid>
      </Container>

      {/* ---------------- 6. MÜŞTERİ YORUMLARI ---------------- */}
      <Container maxWidth="lg" sx={{ py: { xs: 4, md: 6 } }}>
        <Box sx={{ textAlign: 'center', mb: 5 }}>
          {/* Üst Rozet */}
          <Box
            sx={{
              display: 'inline-flex',
              alignItems: 'center',
              gap: 0.75,
              px: 1.75,
              py: 0.5,
              borderRadius: 9999,
              bgcolor: '#EFF6FF',
              border: '1px solid #BFDBFE',
              color: 'primary.main',
              mb: 1.5,
            }}
          >
            <Typography
              variant="caption"
              sx={{
                fontFamily: '"Outfit", sans-serif',
                fontWeight: 700,
                fontSize: '0.8rem',
              }}
            >
              ⭐ Müşteri Memnuniyeti
            </Typography>
          </Box>

          <Typography
            variant="h4"
            component="h2"
            sx={{
              fontFamily: '"Outfit", sans-serif',
              fontWeight: 700,
              color: 'primary.main',
              fontSize: { xs: '1.75rem', md: '2.1rem' },
              letterSpacing: '-0.5px',
              mb: 1,
            }}
          >
            Kullanıcılarımız Ne Diyor?
          </Typography>

          <Typography
            variant="body2"
            sx={{
              fontFamily: '"Inter", sans-serif',
              color: 'text.secondary',
              maxWidth: 500,
              mx: 'auto',
            }}
          >
            Binlerce mutlu müşterimizin mağazamızdan yaptığı alışveriş deneyimlerine göz atın.
          </Typography>
        </Box>

        <Grid container spacing={{ xs: 2.5, sm: 3, md: 3.5 }}>
          {testimonials.map((item) => (
            <Grid size={{ xs: 12, md: 4 }} key={item.id}>
              <Paper
                elevation={0}
                sx={{
                  p: { xs: 3, md: 3.5 },
                  height: '100%',
                  display: 'flex',
                  flexDirection: 'column',
                  borderRadius: 3,
                  bgcolor: '#FFFFFF',
                  border: '1px solid rgba(229, 231, 235, 0.8)',
                  boxShadow: '0 1px 3px 0 rgba(0, 0, 0, 0.05), 0 1px 2px -1px rgba(0, 0, 0, 0.05)',
                  transition: 'all 0.25s ease',
                  '&:hover': {
                    transform: 'translateY(-4px)',
                    boxShadow: '0 10px 20px -3px rgba(30, 58, 138, 0.08)',
                    borderColor: '#BFDBFE',
                  },
                }}
              >
                {/* 5 Yıldız */}
                <Rating
                  value={item.rating}
                  readOnly
                  size="small"
                  sx={{ color: '#F59E0B', mb: 2, fontSize: '1.1rem' }}
                />

                {/* Yorum Metni */}
                <Typography
                  variant="body2"
                  sx={{
                    fontFamily: '"Inter", sans-serif',
                    color: 'text.primary',
                    fontSize: '0.9rem',
                    lineHeight: 1.6,
                    mb: 3,
                    flexGrow: 1,
                    fontStyle: 'italic',
                  }}
                >
                  "{item.comment}"
                </Typography>

                {/* Müşteri Kimlik Çubuğu */}
                <Box sx={{ display: 'flex', alignItems: 'center', gap: 1.5, pt: 2, borderTop: '1px solid #F3F4F6' }}>
                  <Avatar
                    sx={{
                      bgcolor: '#EFF6FF',
                      color: 'primary.main',
                      fontWeight: 700,
                      fontFamily: '"Outfit", sans-serif',
                      fontSize: '0.85rem',
                      width: 42,
                      height: 42,
                      border: '1px solid #BFDBFE',
                    }}
                  >
                    {item.avatar}
                  </Avatar>
                  <Box>
                    <Box sx={{ display: 'flex', alignItems: 'center', gap: 0.5 }}>
                      <Typography
                        variant="subtitle2"
                        sx={{
                          fontFamily: '"Outfit", sans-serif',
                          fontWeight: 700,
                          color: 'primary.main',
                          fontSize: '0.925rem',
                        }}
                      >
                        {item.name}
                      </Typography>
                      <VerifiedIcon sx={{ fontSize: 16, color: 'primary.light' }} />
                    </Box>
                    <Typography
                      variant="caption"
                      sx={{
                        fontFamily: '"Inter", sans-serif',
                        color: 'text.secondary',
                        fontSize: '0.75rem',
                      }}
                    >
                      {item.location}
                    </Typography>
                  </Box>
                </Box>
              </Paper>
            </Grid>
          ))}
        </Grid>
      </Container>

      {/* ---------------- 7. BÜLTEN & 100 TL İNDİRİM KUUPONU ---------------- */}
      <Container maxWidth="lg" sx={{ pt: { xs: 3, md: 5 } }}>
        <Paper
          elevation={0}
          sx={{
            p: { xs: 3, sm: 4, md: 5 },
            borderRadius: 3,
            bgcolor: '#FFFFFF',
            border: '1px solid rgba(229, 231, 235, 0.8)',
            boxShadow: '0 1px 3px 0 rgba(0, 0, 0, 0.05), 0 1px 2px -1px rgba(0, 0, 0, 0.05)',
          }}
        >
          <Grid container spacing={3} alignItems="center" justifyContent="space-between">
            {/* Sol: E-posta İkonu ve Açıklama */}
            <Grid size={{ xs: 12, md: 7 }}>
              <Box sx={{ display: 'flex', alignItems: 'center', gap: 2.5 }}>
                <Box
                  sx={{
                    width: 52,
                    height: 52,
                    borderRadius: 2.5,
                    bgcolor: '#EFF6FF',
                    color: 'primary.main',
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'center',
                    flexShrink: 0,
                    border: '1px solid #BFDBFE',
                  }}
                >
                  <MailOutlineRoundedIcon sx={{ fontSize: 26 }} />
                </Box>
                <Box>
                  <Typography
                    variant="h6"
                    sx={{
                      fontFamily: '"Outfit", sans-serif',
                      fontWeight: 700,
                      color: 'primary.main',
                      fontSize: { xs: '1.05rem', sm: '1.25rem' },
                    }}
                  >
                    İlk Siparişinize Özel 100 TL İndirim
                  </Typography>
                  <Typography
                    variant="body2"
                    sx={{
                      fontFamily: '"Inter", sans-serif',
                      color: 'text.secondary',
                      fontSize: '0.85rem',
                      mt: 0.3,
                    }}
                  >
                    Bültenimize kaydolun, yeni ürün ve özel fırsatlardan anında haberdar olun.
                  </Typography>
                </Box>
              </Box>
            </Grid>

            {/* Sağ: Form Girdisi ve Buton */}
            <Grid size={{ xs: 12, md: 5 }}>
              <Box
                component="form"
                onSubmit={handleNewsletterSubmit}
                sx={{
                  display: 'flex',
                  gap: 1.5,
                  alignItems: 'center',
                  width: '100%',
                }}
              >
                <OutlinedInput
                  fullWidth
                  placeholder="E-posta adresinizi girin"
                  value={newsletterEmail}
                  onChange={(e) => setNewsletterEmail(e.target.value)}
                  sx={{
                    height: 44,
                    borderRadius: '10px',
                    bgcolor: '#FFFFFF',
                    fontSize: '0.875rem',
                    fontFamily: '"Inter", sans-serif',
                    '& fieldset': { borderColor: '#D1D5DB' },
                    '&:hover fieldset': { borderColor: '#9CA3AF' },
                    '&.Mui-focused fieldset': { borderColor: 'primary.light' },
                  }}
                />
                <Button
                  type="submit"
                  variant="contained"
                  sx={{
                    height: 44,
                    px: 3.5,
                    borderRadius: '10px',
                    bgcolor: 'primary.main',
                    color: '#FFFFFF',
                    fontWeight: 700,
                    fontSize: '0.9rem',
                    flexShrink: 0,
                    '&:hover': {
                      bgcolor: 'primary.dark',
                    },
                  }}
                >
                  Kayıt Ol
                </Button>
              </Box>
            </Grid>
          </Grid>
        </Paper>
      </Container>
    </Box>
  );
}
