import { useEffect, useState } from 'react';
import { useNavigate, useParams, Link as RouterLink } from 'react-router';
import {
  Box,
  Container,
  Typography,
  Stack,
  Button,
  IconButton,
  Paper,
  Grid,
  Breadcrumbs,
  Link as MuiLink,
  Tooltip,
} from '@mui/material';
import ArrowBackIcon from '@mui/icons-material/ArrowBack';
import ShoppingCartIcon from '@mui/icons-material/ShoppingCart';
import ShoppingCartOutlinedIcon from '@mui/icons-material/ShoppingCartOutlined';
import CheckCircleOutlineIcon from '@mui/icons-material/CheckCircleOutline';
import StarIcon from '@mui/icons-material/Star';
import FavoriteBorderIcon from '@mui/icons-material/FavoriteBorder';
import FavoriteIcon from '@mui/icons-material/Favorite';
import ShareOutlinedIcon from '@mui/icons-material/ShareOutlined';
import LocalShippingOutlinedIcon from '@mui/icons-material/LocalShippingOutlined';
import VerifiedUserOutlinedIcon from '@mui/icons-material/VerifiedUserOutlined';
import CachedOutlinedIcon from '@mui/icons-material/CachedOutlined';
import SecurityOutlinedIcon from '@mui/icons-material/SecurityOutlined';
import AssignmentReturnOutlinedIcon from '@mui/icons-material/AssignmentReturnOutlined';
import SupportAgentOutlinedIcon from '@mui/icons-material/SupportAgentOutlined';
import KeyboardArrowLeftIcon from '@mui/icons-material/KeyboardArrowLeft';
import KeyboardArrowRightIcon from '@mui/icons-material/KeyboardArrowRight';
import AddIcon from '@mui/icons-material/Add';
import RemoveIcon from '@mui/icons-material/Remove';
import MenuBookOutlinedIcon from '@mui/icons-material/MenuBookOutlined';
import CreateOutlinedIcon from '@mui/icons-material/CreateOutlined';
import DescriptionOutlinedIcon from '@mui/icons-material/DescriptionOutlined';
import ViewDayOutlinedIcon from '@mui/icons-material/ViewDayOutlined';
import LayersOutlinedIcon from '@mui/icons-material/LayersOutlined';
import AspectRatioOutlinedIcon from '@mui/icons-material/AspectRatioOutlined';
import CheckIcon from '@mui/icons-material/Check';
import NavigateNextIcon from '@mui/icons-material/NavigateNext';
import { toast } from 'react-toastify';

import { productService } from '../Services/productService';
import type { FullProductresponse, ProductDetailProps, GetProducts } from '../models/ProductResponse';
import Spinner from './Spinner';

export default function ProductDetail({ addtocart, cartItems }: ProductDetailProps) {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();

  const [product, setProduct] = useState<FullProductresponse | null>(null);
  const [similarProducts, setSimilarProducts] = useState<GetProducts[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [selectedImageIndex, setSelectedImageIndex] = useState<number>(0);
  const [quantity, setQuantity] = useState<number>(1);
  const [isFavorited, setIsFavorited] = useState<boolean>(false);
  const [favoriteSimilarIds, setFavoriteSimilarIds] = useState<Record<string, boolean>>({});

  useEffect(() => {
    let isMounted = true;
    const loadProductData = async () => {
      if (!id) return;
      try {
        setLoading(true);
        window.scrollTo({ top: 0, behavior: 'smooth' });

        const productData = await productService.getProductById(id);
        if (isMounted) {
          setProduct(productData);
          setSelectedImageIndex(0);
          setQuantity(1);
        }

        // Benzer ürünleri API'den çek
        const simRes = await productService.getPublicProducts(0, 10);
        if (isMounted && simRes && simRes.Data) {
          const filtered = simRes.Data.filter((p) => p.product_id !== id).slice(0, 4);
          setSimilarProducts(filtered);
        }
      } catch (error) {
        console.error('Ürün detayları alınırken hata oluştu:', error);
      } finally {
        if (isMounted) setLoading(false);
      }
    };

    loadProductData();
    return () => {
      isMounted = false;
    };
  }, [id]);

  if (loading) {
    return <Spinner message="Ürün detayları yükleniyor..." />;
  }

  if (!product) {
    return (
      <Box sx={{ py: 12, textAlign: 'center', bgcolor: '#F8FAFC', minHeight: '60vh' }}>
        <Container maxWidth="sm">
          <Typography variant="h4" fontWeight="bold" color="text.primary" gutterBottom>
            Ürün Bulunamadı
          </Typography>
          <Typography variant="body1" color="text.secondary" sx={{ mb: 4 }}>
            Aradığınız ürün mevcut değil veya kaldırılmış olabilir.
          </Typography>
          <Button
            variant="contained"
            color="primary"
            onClick={() => navigate('/products')}
            sx={{ px: 4, py: 1.5, borderRadius: 2 }}
          >
            Ürünlere Göz At
          </Button>
        </Container>
      </Box>
    );
  }

  const isInCart = cartItems ? cartItems.some((item) => item.product_id === product.ProductID) : false;

  // Görseller listesi (API'den boş gelirse fallback oluştur)
  const productImages =
    product.Images && product.Images.length > 0
      ? product.Images
      : [{ image_id: 'default', path: 'https://images.unsplash.com/photo-1544816155-12df9643f363?w=800' }];

  const currentImage = productImages[selectedImageIndex] || productImages[0];

  const handlePrevImage = () => {
    setSelectedImageIndex((prev) => (prev > 0 ? prev - 1 : productImages.length - 1));
  };

  const handleNextImage = () => {
    setSelectedImageIndex((prev) => (prev < productImages.length - 1 ? prev + 1 : 0));
  };

  const handleQuantityMinus = () => {
    setQuantity((prev) => (prev > 1 ? prev - 1 : 1));
  };

  const handleQuantityPlus = () => {
    if (product.Stock) {
          const maxStock = product.Stock > 0 ? product.Stock : 99;
    setQuantity((prev) => (prev < maxStock ? prev + 1 : prev));
    }

  };

  const handleAddToCart = () => {
    addtocart({
      product_id: product.ProductID ?? '',
      product_name: product.ProductName ?? '',
      price: product.Price ?? 0,
      stock: product.Stock ?? 0,
      thumbnail_image: currentImage.path,
      quantity: quantity,
    });
  };

  const handleToggleFavorite = () => {
    setIsFavorited((prev) => !prev);
    if (!isFavorited) {
      toast.success('Ürün favorilere eklendi!');
    } else {
      toast.info('Ürün favorilerden çıkarıldı.');
    }
  };

  const handleToggleSimilarFavorite = (simId: string) => {
    setFavoriteSimilarIds((prev) => ({ ...prev, [simId]: !prev[simId] }));
  };

  const handleShare = () => {
    if (navigator.clipboard) {
      navigator.clipboard.writeText(window.location.href);
      toast.success('Ürün bağlantısı panoya kopyalandı!');
    }
  };

  // İndirim ve eski fiyat hesaplama
  const discountedPrice = product.Price ?? 0;
  const originalPrice = Math.round(discountedPrice * 1.17);
  const discountPercent = Math.round(((originalPrice - discountedPrice) / originalPrice) * 100);

  // Kategori adı
  const categoryName =
    product.Categories && product.Categories.length > 0
      ? product.Categories[0].category_name
      : 'Kırtasiye & Ofis';

  // Teknik Detaylar Verisi
  // API'den gelen özellikler varsa onları kullan, yoksa tasarımdaki zengin örnekleri göster
  const defaultSpecs = [
    { label: 'Ciltli Kapak', value: 'Evet (Sert Karton Kaplama)', icon: <MenuBookOutlinedIcon fontSize="small" /> },
    { label: 'Hediye Kalem', value: 'Evet (Mat Jel Kalem Dahil)', isHighlight: true, icon: <CreateOutlinedIcon fontSize="small" /> },
    { label: 'Sayfa Sayısı', value: '64 - 84 - 100 Sayfa', icon: <DescriptionOutlinedIcon fontSize="small" /> },
    { label: 'Spiral / Telli', value: 'Evet (1 Adet Çift Spiralli)', icon: <ViewDayOutlinedIcon fontSize="small" /> },
    { label: 'Kağıt Gramajı', value: '80 gr / Ivory Fildişi', icon: <LayersOutlinedIcon fontSize="small" /> },
    { label: 'Format & Ebat', value: 'A5 & B5 Seçenekleri', icon: <AspectRatioOutlinedIcon fontSize="small" /> },
  ];

  const technicalDetailsList =
    product.TechnicalDetails && Object.keys(product.TechnicalDetails).length > 0
      ? Object.entries(product.TechnicalDetails).map(([key, value], idx) => {
          const specIcons = [
            <MenuBookOutlinedIcon fontSize="small" key="1" />,
            <CreateOutlinedIcon fontSize="small" key="2" />,
            <DescriptionOutlinedIcon fontSize="small" key="3" />,
            <ViewDayOutlinedIcon fontSize="small" key="4" />,
            <LayersOutlinedIcon fontSize="small" key="5" />,
            <AspectRatioOutlinedIcon fontSize="small" key="6" />,
          ];
          return {
            label: key,
            value: String(value),
            isHighlight: idx === 1,
            icon: specIcons[idx % specIcons.length],
          };
        })
      : defaultSpecs;

  // Benzer ürünler için sabit örnek rozetler ve değerlendirmeler (Tasarım referansı)
  const similarCardMeta = [
    { badge: 'YENİ', badgeColor: '#2563EB', tag: 'Ücretsiz Kargo', rating: '5.0', reviews: 34, oldPriceMult: 1.12, category: 'GİYİM & MODA' },
    { badge: 'FIRSAT', badgeColor: '#EA580C', tag: 'Hızlı Gönderi', rating: '4.9', reviews: 116, oldPriceMult: 1.13, category: 'ÖZEL SERİ' },
    { badge: '%20 İNDİRİM', badgeColor: '#EA580C', tag: 'Popüler', rating: '4.8', reviews: 58, oldPriceMult: 1.25, category: 'GİYİM & MODA' },
    { badge: 'TÜKENİYOR', badgeColor: '#0F275A', tag: 'Stokta Az', rating: '4.9', reviews: 162, oldPriceMult: 1.2, category: 'KIRTASİYE & OFİS' },
  ];

  return (
    <Box sx={{ bgcolor: '#F8FAFC', minHeight: '100vh', py: { xs: 2, md: 4 } }}>
      <Container maxWidth="lg">
        {/* Üst Gezinme: Geri Dön Butonu ve Breadcrumbs */}
        <Stack
          direction="row"
          justifyContent="space-between"
          alignItems="center"
          sx={{ mb: 3 }}
        >
          <Button
            startIcon={<ArrowBackIcon fontSize="small" />}
            onClick={() => navigate(-1)}
            sx={{
              textTransform: 'none',
              fontWeight: 600,
              fontSize: '0.925rem',
              color: '#334155',
              px: 1.5,
              py: 0.5,
              borderRadius: 2,
              '&:hover': { bgcolor: '#E2E8F0', color: '#0F172A' },
            }}
          >
            Geri Dön
          </Button>

          <Breadcrumbs
            separator={<NavigateNextIcon fontSize="small" sx={{ color: '#94A3B8' }} />}
            aria-label="breadcrumb"
            sx={{ fontSize: '0.85rem' }}
          >
            <MuiLink
              component={RouterLink}
              to="/"
              underline="hover"
              color="#64748B"
              sx={{ fontWeight: 500 }}
            >
              Ana Sayfa
            </MuiLink>
            <MuiLink
              component={RouterLink}
              to="/products"
              underline="hover"
              color="#64748B"
              sx={{ fontWeight: 500 }}
            >
              {categoryName}
            </MuiLink>
            <Typography color="#0F172A" sx={{ fontWeight: 700, maxWidth: { xs: 150, sm: 'none' } }} noWrap>
              {product.ProductName}
            </Typography>
          </Breadcrumbs>
        </Stack>

        {/* 1. ANA ÜRÜN KARTI (2 Kolonlu) */}
        <Paper
          elevation={0}
          sx={{
            p: { xs: 2.5, md: 4 },
            borderRadius: 3.5,
            bgcolor: '#FFFFFF',
            border: '1px solid #E2E8F0',
            boxShadow: '0 4px 20px -2px rgba(0, 0, 0, 0.03)',
          }}
        >
          <Grid container spacing={{ xs: 3, md: 5 }}>
            {/* SOL KOLON: Ürün Görseli, Etiketler, Oklar ve Küçük Resimler */}
            <Grid size={{ xs: 12, md: 6 }}>
              {/* Ana Görsel Kutusu */}
              <Box
                sx={{
                  position: 'relative',
                  width: '100%',
                  height: { xs: 340, sm: 420, md: 470 },
                  bgcolor: '#F8FAFC',
                  borderRadius: 3,
                  border: '1px solid #E2E8F0',
                  overflow: 'hidden',
                  display: 'flex',
                  alignItems: 'center',
                  justifyContent: 'center',
                }}
              >
                {/* Sol Üstteki Rozetler */}
                <Stack
                  direction="column"
                  spacing={0.75}
                  sx={{
                    position: 'absolute',
                    top: 16,
                    left: 16,
                    zIndex: 2,
                  }}
                >
                  <Box
                    sx={{
                      bgcolor: '#0F275A',
                      color: '#FFFFFF',
                      fontSize: '0.72rem',
                      fontWeight: 800,
                      letterSpacing: '0.5px',
                      px: 1.5,
                      py: 0.5,
                      borderRadius: 1.5,
                      textTransform: 'uppercase',
                      boxShadow: '0 2px 8px rgba(15, 39, 90, 0.25)',
                    }}
                  >
                    EN ÇOK SATAN
                  </Box>
                  <Box
                    sx={{
                      bgcolor: '#F97316',
                      color: '#FFFFFF',
                      fontSize: '0.72rem',
                      fontWeight: 800,
                      letterSpacing: '0.5px',
                      px: 1.5,
                      py: 0.5,
                      borderRadius: 1.5,
                      textTransform: 'uppercase',
                      boxShadow: '0 2px 8px rgba(249, 115, 22, 0.25)',
                      width: 'fit-content',
                    }}
                  >
                    %{discountPercent} İNDİRİM
                  </Box>
                </Stack>

                {/* Sol Ok Butonu */}
                <IconButton
                  onClick={handlePrevImage}
                  sx={{
                    position: 'absolute',
                    left: 14,
                    top: '50%',
                    transform: 'translateY(-50%)',
                    zIndex: 2,
                    bgcolor: 'rgba(255, 255, 255, 0.95)',
                    boxShadow: '0 2px 10px rgba(0, 0, 0, 0.12)',
                    color: '#334155',
                    width: 38,
                    height: 38,
                    '&:hover': { bgcolor: '#FFFFFF', color: '#0F172A' },
                  }}
                  aria-label="Önceki görsel"
                >
                  <KeyboardArrowLeftIcon />
                </IconButton>

                {/* Sağ Ok Butonu */}
                <IconButton
                  onClick={handleNextImage}
                  sx={{
                    position: 'absolute',
                    right: 14,
                    top: '50%',
                    transform: 'translateY(-50%)',
                    zIndex: 2,
                    bgcolor: 'rgba(255, 255, 255, 0.95)',
                    boxShadow: '0 2px 10px rgba(0, 0, 0, 0.12)',
                    color: '#334155',
                    width: 38,
                    height: 38,
                    '&:hover': { bgcolor: '#FFFFFF', color: '#0F172A' },
                  }}
                  aria-label="Sonraki görsel"
                >
                  <KeyboardArrowRightIcon />
                </IconButton>

                {/* Ana Resim */}
                <Box
                  component="img"
                  src={currentImage.path}
                  alt={product.ProductName}
                  sx={{
                    maxWidth: '100%',
                    maxHeight: '100%',
                    objectFit: 'contain',
                    transition: 'all 0.3s ease',
                  }}
                />
              </Box>

              {/* Alt Thumbnail Listesi ve Sayfalama Noktaları */}
              <Box
                sx={{
                  display: 'flex',
                  alignItems: 'center',
                  justifyContent: 'space-between',
                  mt: 2.5,
                  px: 0.5,
                }}
              >
                {/* Thumbnails */}
                <Stack direction="row" spacing={1.5}>
                  {productImages.slice(0, 4).map((img, idx) => (
                    <Box
                      key={img.image_id || idx}
                      onClick={() => setSelectedImageIndex(idx)}
                      sx={{
                        width: 64,
                        height: 64,
                        borderRadius: 2,
                        border: '2px solid',
                        borderColor: selectedImageIndex === idx ? '#2563EB' : '#E2E8F0',
                        overflow: 'hidden',
                        cursor: 'pointer',
                        bgcolor: '#FFFFFF',
                        p: 0.5,
                        transition: 'all 0.2s ease',
                        boxShadow:
                          selectedImageIndex === idx ? '0 0 0 2px rgba(37, 99, 235, 0.2)' : 'none',
                        '&:hover': {
                          borderColor: '#2563EB',
                        },
                      }}
                    >
                      <Box
                        component="img"
                        src={img.path}
                        alt={`thumbnail-${idx}`}
                        sx={{
                          width: '100%',
                          height: '100%',
                          objectFit: 'cover',
                          borderRadius: 1,
                        }}
                      />
                    </Box>
                  ))}
                </Stack>

                {/* Sayfalama Noktaları (Dots) */}
                <Stack direction="row" spacing={0.8} alignItems="center">
                  {productImages.map((_, idx) => (
                    <Box
                      key={idx}
                      onClick={() => setSelectedImageIndex(idx)}
                      sx={{
                        cursor: 'pointer',
                        width: selectedImageIndex === idx ? 16 : 8,
                        height: 8,
                        borderRadius: 4,
                        bgcolor: selectedImageIndex === idx ? '#0F275A' : '#CBD5E1',
                        transition: 'all 0.3s cubic-bezier(0.4, 0, 0.2, 1)',
                      }}
                    />
                  ))}
                </Stack>
              </Box>
            </Grid>

            {/* SAĞ KOLON: Ürün Detayları, Fiyat, Satın Alma ve Rozetler */}
            <Grid size={{ xs: 12, md: 6 }}>
              <Stack spacing={2.5}>
                {/* Üst Etiket & Puan */}
                <Stack direction="row" justifyContent="space-between" alignItems="center">
                  <Typography
                    variant="caption"
                    sx={{
                      color: '#2563EB',
                      fontWeight: 800,
                      letterSpacing: '0.8px',
                      textTransform: 'uppercase',
                      fontSize: '0.8rem',
                    }}
                  >
                    PREMİUM KOLEKSİYON
                  </Typography>

                  <Stack direction="row" spacing={0.6} alignItems="center">
                    <StarIcon sx={{ color: '#F59E0B', fontSize: '1.2rem' }} />
                    <Typography variant="body2" sx={{ fontWeight: 800, color: '#0F172A' }}>
                      4.9
                    </Typography>
                    <Typography variant="body2" sx={{ color: '#64748B' }}>
                      (128 Yorum)
                    </Typography>
                  </Stack>
                </Stack>

                {/* Ürün Başlığı */}
                <Typography
                  variant="h1"
                  sx={{
                    fontFamily: '"Outfit", sans-serif',
                    fontWeight: 800,
                    fontSize: { xs: '1.75rem', sm: '2.25rem', md: '2.4rem' },
                    color: '#0F172A',
                    letterSpacing: '-0.5px',
                    lineHeight: 1.2,
                  }}
                >
                  {product.ProductName}
                </Typography>

                {/* Stok & Kargo Durumu */}
                <Stack direction="row" spacing={0.8} alignItems="center">
                  <CheckCircleOutlineIcon sx={{ color: '#16A34A', fontSize: '1.15rem' }} />
                  <Typography variant="body2" sx={{ color: '#16A34A', fontWeight: 700 }}>
                    {product.Stock ?? 0 > 0 ? 'Stokta Var • 24 Saatte Kargoda' : 'Tükendi'}
                  </Typography>
                </Stack>

                {/* Fiyat Kutusu */}
                <Box
                  sx={{
                    bgcolor: '#F0F7FF',
                    borderRadius: 3,
                    p: { xs: 2, sm: 2.5 },
                    border: '1px solid #E0F2FE',
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'space-between',
                  }}
                >
                  <Stack direction="row" spacing={2} alignItems="baseline">
                    <Typography
                      variant="h3"
                      sx={{
                        fontFamily: '"Outfit", sans-serif',
                        fontWeight: 800,
                        fontSize: { xs: '2rem', sm: '2.4rem' },
                        color: '#0F275A',
                        letterSpacing: '-0.5px',
                      }}
                    >
                      {discountedPrice.toLocaleString('tr-TR')} ₺
                    </Typography>
                    <Typography
                      variant="body1"
                      sx={{
                        color: '#94A3B8',
                        textDecoration: 'line-through',
                        fontSize: '1.1rem',
                        fontWeight: 500,
                      }}
                    >
                      {originalPrice.toLocaleString('tr-TR')} ₺
                    </Typography>
                  </Stack>

                  <Box
                    sx={{
                      bgcolor: '#FFEDD5',
                      color: '#EA580C',
                      fontWeight: 800,
                      fontSize: '0.85rem',
                      px: 2,
                      py: 0.7,
                      borderRadius: 9999,
                      boxShadow: '0 1px 3px rgba(234, 88, 12, 0.1)',
                    }}
                  >
                    %{discountPercent} Kazanç
                  </Box>
                </Box>

                {/* Ürün Açıklaması */}
                <Box sx={{ pt: 0.5 }}>
                  <Typography
                    variant="subtitle1"
                    sx={{
                      fontFamily: '"Outfit", sans-serif',
                      fontWeight: 800,
                      color: '#0F172A',
                      fontSize: '1.05rem',
                      mb: 1,
                    }}
                  >
                    Ürün Açıklaması
                  </Typography>
                  <Typography
                    variant="body2"
                    sx={{
                      color: '#475569',
                      lineHeight: 1.75,
                      fontSize: '0.925rem',
                    }}
                  >
                    {product.Description ||
                      'Not Defteri Seti: 1 Ciltli Klasik Defter, 1 Spiral Telli Not Defteri, 1 Büyük Boy Eskiz Defteri ve 1 adet Minimalist Roller Kalem hediyelidir. Günlük ajanda, toplantı notları ve yaratıcı fikirleriniz için fildişi dokulu kağıt kalitesi sunar.'}
                  </Typography>
                </Box>

                {/* Adet Seçici & Sepete Ekle Butonu */}
                <Stack direction="row" spacing={2} sx={{ pt: 1 }}>
                  {/* Stepper (Miktar Artırma/Azaltma) */}
                  <Box
                    sx={{
                      display: 'flex',
                      alignItems: 'center',
                      border: '1.5px solid #E2E8F0',
                      borderRadius: 2.5,
                      bgcolor: '#FFFFFF',
                      px: 0.5,
                      py: 0.25,
                    }}
                  >
                    <IconButton
                      size="small"
                      onClick={handleQuantityMinus}
                      disabled={quantity <= 1}
                      sx={{ color: '#475569' }}
                    >
                      <RemoveIcon fontSize="small" />
                    </IconButton>
                    <Typography
                      sx={{
                        minWidth: 36,
                        textAlign: 'center',
                        fontWeight: 700,
                        fontSize: '1rem',
                        color: '#0F172A',
                      }}
                    >
                      {quantity}
                    </Typography>
                    <IconButton
                      size="small"
                      onClick={handleQuantityPlus}
                      disabled={(product.Stock ?? 0) <= 0 || quantity >= (product.Stock ?? 0)}
                      sx={{ color: '#475569' }}
                    >
                      <AddIcon fontSize="small" />
                    </IconButton>
                  </Box>

                  {/* Sepete Ekle Butonu */}
                  <Button
                    variant="contained"
                    fullWidth
                    startIcon={<ShoppingCartIcon />}
                    onClick={handleAddToCart}
                    disabled={(product.Stock ?? 0) <= 0}
                    sx={{
                      bgcolor: '#F97316',
                      color: '#FFFFFF',
                      py: 1.5,
                      px: 3,
                      borderRadius: 2.5,
                      fontSize: '1.025rem',
                      fontWeight: 700,
                      textTransform: 'none',
                      boxShadow: '0 4px 14px 0 rgba(249, 115, 22, 0.35)',
                      '&:hover': {
                        bgcolor: '#EA580C',
                        boxShadow: '0 6px 18px 0 rgba(234, 88, 12, 0.45)',
                      },
                      '&.Mui-disabled': {
                        bgcolor: '#CBD5E1',
                        color: '#94A3B8',
                      },
                    }}
                  >
                    {isInCart ? 'Sepete Tekrar Ekle' : 'Sepete Ekle'}
                  </Button>
                </Stack>

                {/* İkincil Eylemler: Favorilere Ekle & Paylaş */}
                <Stack direction="row" spacing={2}>
                  <Button
                    variant="text"
                    startIcon={
                      isFavorited ? (
                        <FavoriteIcon sx={{ color: '#EF4444' }} fontSize="small" />
                      ) : (
                        <FavoriteBorderIcon fontSize="small" />
                      )
                    }
                    onClick={handleToggleFavorite}
                    sx={{
                      flex: 1,
                      py: 1.2,
                      borderRadius: 2,
                      bgcolor: isFavorited ? '#FEF2F2' : '#F1F5F9',
                      color: isFavorited ? '#EF4444' : '#334155',
                      fontWeight: 600,
                      fontSize: '0.9rem',
                      textTransform: 'none',
                      '&:hover': {
                        bgcolor: isFavorited ? '#FEE2E2' : '#E2E8F0',
                      },
                    }}
                  >
                    {isFavorited ? 'Favorilerden Çıkar' : 'Favorilere Ekle'}
                  </Button>

                  <Button
                    variant="text"
                    startIcon={<ShareOutlinedIcon fontSize="small" />}
                    onClick={handleShare}
                    sx={{
                      flex: 1,
                      py: 1.2,
                      borderRadius: 2,
                      bgcolor: '#F1F5F9',
                      color: '#334155',
                      fontWeight: 600,
                      fontSize: '0.9rem',
                      textTransform: 'none',
                      '&:hover': { bgcolor: '#E2E8F0' },
                    }}
                  >
                    Paylaş
                  </Button>
                </Stack>

                {/* 3'lü Güven & Teslimat Rozetleri */}
                <Box
                  sx={{
                    bgcolor: '#F8FAFC',
                    borderRadius: 2.5,
                    p: 2,
                    border: '1px solid #EEF2F6',
                    mt: 1,
                  }}
                >
                  <Grid container spacing={2} justifyContent="space-between">
                    <Grid size={{ xs: 4 }} sx={{ textAlign: 'center' }}>
                      <LocalShippingOutlinedIcon sx={{ color: '#1E3A8A', fontSize: 26, mb: 0.5 }} />
                      <Typography variant="body2" sx={{ fontWeight: 800, color: '#0F172A', fontSize: '0.85rem' }}>
                        Ücretsiz Kargo
                      </Typography>
                      <Typography variant="caption" sx={{ color: '#64748B', display: 'block', fontSize: '0.75rem' }}>
                        Tüm siparişlerde
                      </Typography>
                    </Grid>

                    <Grid size={{ xs: 4 }} sx={{ textAlign: 'center', borderLeft: '1px solid #E2E8F0', borderRight: '1px solid #E2E8F0' }}>
                      <VerifiedUserOutlinedIcon sx={{ color: '#1E3A8A', fontSize: 26, mb: 0.5 }} />
                      <Typography variant="body2" sx={{ fontWeight: 800, color: '#0F172A', fontSize: '0.85rem' }}>
                        2 Yıl Garanti
                      </Typography>
                      <Typography variant="caption" sx={{ color: '#64748B', display: 'block', fontSize: '0.75rem' }}>
                        Orijinal ürün
                      </Typography>
                    </Grid>

                    <Grid size={{ xs: 4 }} sx={{ textAlign: 'center' }}>
                      <CachedOutlinedIcon sx={{ color: '#1E3A8A', fontSize: 26, mb: 0.5 }} />
                      <Typography variant="body2" sx={{ fontWeight: 800, color: '#0F172A', fontSize: '0.85rem' }}>
                        14 Gün İade
                      </Typography>
                      <Typography variant="caption" sx={{ color: '#64748B', display: 'block', fontSize: '0.75rem' }}>
                        Koşulsuz güvence
                      </Typography>
                    </Grid>
                  </Grid>
                </Box>
              </Stack>
            </Grid>
          </Grid>
        </Paper>

        {/* 2. TEKNİK DETAYLAR KARTI */}
        <Paper
          elevation={0}
          sx={{
            p: { xs: 3, md: 4 },
            mt: 4,
            borderRadius: 3.5,
            bgcolor: '#FFFFFF',
            border: '1px solid #E2E8F0',
          }}
        >
          {/* Başlık Satırı */}
          <Stack direction="row" justifyContent="space-between" alignItems="center" sx={{ mb: 3 }}>
            <Stack direction="row" spacing={1.5} alignItems="center">
              <Box
                sx={{
                  width: 5,
                  height: 24,
                  bgcolor: '#0F275A',
                  borderRadius: 1,
                }}
              />
              <Typography
                variant="h5"
                sx={{
                  fontFamily: '"Outfit", sans-serif',
                  fontWeight: 800,
                  color: '#0F172A',
                  fontSize: { xs: '1.25rem', md: '1.45rem' },
                }}
              >
                Teknik Detaylar
              </Typography>
            </Stack>

            <Typography
              variant="caption"
              sx={{
                color: '#94A3B8',
                fontWeight: 700,
                letterSpacing: '1px',
                fontSize: '0.78rem',
              }}
            >
              ÖZELLİK LİSTESİ
            </Typography>
          </Stack>

          {/* Özellik Kutuları Izgarası */}
          <Grid container spacing={2}>
            {technicalDetailsList.map((spec, index) => (
              <Grid size={{ xs: 12, md: 6 }} key={index}>
                <Box
                  sx={{
                    bgcolor: '#F8FAFC',
                    border: '1px solid #EEF2F6',
                    borderRadius: 2.5,
                    p: 2,
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'space-between',
                  }}
                >
                  <Stack direction="row" spacing={1.5} alignItems="center">
                    <Box sx={{ color: '#64748B', display: 'flex', alignItems: 'center' }}>
                      {spec.icon}
                    </Box>
                    <Typography
                      variant="body2"
                      sx={{
                        color: '#334155',
                        fontWeight: 600,
                        fontSize: '0.9rem',
                      }}
                    >
                      {spec.label}
                    </Typography>
                  </Stack>

                  <Typography
                    variant="body2"
                    sx={{
                      fontWeight: 700,
                      fontSize: '0.9rem',
                      color: spec.isHighlight ? '#EA580C' : '#0F172A',
                    }}
                  >
                    {spec.value}
                  </Typography>
                </Box>
              </Grid>
            ))}
          </Grid>
        </Paper>

        {/* 3. BENZER ÜRÜNLER BÖLÜMÜ */}
        {similarProducts.length > 0 && (
          <Box sx={{ mt: 6 }}>
            {/* Başlık & "Tümünü Gör >" */}
            <Stack
              direction="row"
              justifyContent="space-between"
              alignItems="flex-end"
              sx={{ mb: 3 }}
            >
              <Box>
                <Typography
                  variant="h4"
                  sx={{
                    fontFamily: '"Outfit", sans-serif',
                    fontWeight: 800,
                    color: '#0F172A',
                    fontSize: { xs: '1.5rem', md: '1.85rem' },
                  }}
                >
                  Benzer Ürünler
                </Typography>
                <Typography variant="body2" sx={{ color: '#64748B', mt: 0.5 }}>
                  Bu ürünü inceleyenler bu seçenekleri de tercih etti.
                </Typography>
              </Box>

              <MuiLink
                component={RouterLink}
                to="/products"
                underline="none"
                sx={{
                  color: '#0F275A',
                  fontWeight: 700,
                  fontSize: '0.9rem',
                  display: 'flex',
                  alignItems: 'center',
                  '&:hover': { color: '#2563EB', textDecoration: 'underline' },
                }}
              >
                Tümünü Gör &gt;
              </MuiLink>
            </Stack>

            {/* 4'lü Ürün Kartı Izgarası */}
            <Grid container spacing={3}>
              {similarProducts.map((simProd, idx) => {
                const meta = similarCardMeta[idx % similarCardMeta.length];
                const simInCart = cartItems ? cartItems.some((i) => i.product_id === simProd.product_id) : false;
                const isSimFav = !!favoriteSimilarIds[simProd.product_id];
                const simOldPrice = Math.round(simProd.price * meta.oldPriceMult);

                return (
                  <Grid size={{ xs: 12, sm: 6, md: 3 }} key={simProd.product_id}>
                    <Paper
                      elevation={0}
                      onClick={() => navigate(`/product/${simProd.product_id}`)}
                      sx={{
                        borderRadius: 3.5,
                        overflow: 'hidden',
                        border: '1px solid #E2E8F0',
                        bgcolor: '#FFFFFF',
                        height: '100%',
                        display: 'flex',
                        flexDirection: 'column',
                        cursor: 'pointer',
                        transition: 'all 0.3s cubic-bezier(0.4, 0, 0.2, 1)',
                        '&:hover': {
                          transform: 'translateY(-6px)',
                          boxShadow: '0 16px 30px rgba(15, 39, 90, 0.08)',
                        },
                      }}
                    >
                      {/* Görsel Alanı + Rozet + Favori Kalp */}
                      <Box
                        sx={{
                          position: 'relative',
                          height: 230,
                          bgcolor: '#FAFAFA',
                          display: 'flex',
                          alignItems: 'center',
                          justifyContent: 'center',
                          p: 2,
                          overflow: 'hidden',
                        }}
                      >
                        {/* Sol Üst Rozet */}
                        <Box
                          sx={{
                            position: 'absolute',
                            top: 12,
                            left: 12,
                            bgcolor: meta.badgeColor,
                            color: '#FFFFFF',
                            fontSize: '0.68rem',
                            fontWeight: 800,
                            letterSpacing: '0.5px',
                            px: 1.2,
                            py: 0.4,
                            borderRadius: 1.2,
                            zIndex: 1,
                          }}
                        >
                          {meta.badge}
                        </Box>

                        {/* Sağ Üst Favori Butonu */}
                        <Tooltip title="Favorilere Ekle">
                          <IconButton
                            size="small"
                            onClick={(e) => {
                              e.stopPropagation();
                              handleToggleSimilarFavorite(simProd.product_id);
                            }}
                            sx={{
                              position: 'absolute',
                              top: 12,
                              right: 12,
                              bgcolor: 'rgba(255, 255, 255, 0.9)',
                              boxShadow: '0 2px 6px rgba(0, 0, 0, 0.08)',
                              color: isSimFav ? '#EF4444' : '#64748B',
                              zIndex: 1,
                              '&:hover': { bgcolor: '#FFFFFF', color: '#EF4444' },
                            }}
                          >
                            {isSimFav ? <FavoriteIcon fontSize="small" /> : <FavoriteBorderIcon fontSize="small" />}
                          </IconButton>
                        </Tooltip>

                        {/* Ürün Görseli */}
                        <Box
                          component="img"
                          src={simProd.thumbnail_image}
                          alt={simProd.product_name}
                          sx={{
                            maxHeight: '100%',
                            maxWidth: '100%',
                            objectFit: 'contain',
                            transition: 'transform 0.3s ease',
                            '&:hover': { transform: 'scale(1.05)' },
                          }}
                        />
                      </Box>

                      {/* Kart Bilgi Alanı */}
                      <Box sx={{ p: 2.2, display: 'flex', flexDirection: 'column', flexGrow: 1 }}>
                        <Typography
                          variant="caption"
                          sx={{
                            color: '#64748B',
                            fontWeight: 700,
                            letterSpacing: '0.5px',
                            textTransform: 'uppercase',
                            fontSize: '0.72rem',
                          }}
                        >
                          {meta.category}
                        </Typography>

                        <Typography
                          variant="subtitle2"
                          noWrap
                          sx={{
                            fontFamily: '"Outfit", sans-serif',
                            fontWeight: 700,
                            color: '#0F172A',
                            fontSize: '0.95rem',
                            mt: 0.5,
                            mb: 0.5,
                          }}
                        >
                          {simProd.product_name}
                        </Typography>

                        {/* Yıldız Değerlendirmesi */}
                        <Stack direction="row" spacing={0.5} alignItems="center" sx={{ mb: 1.5 }}>
                          <Stack direction="row" spacing={0.2}>
                            {[...Array(5)].map((_, i) => (
                              <StarIcon key={i} sx={{ color: '#F59E0B', fontSize: '0.9rem' }} />
                            ))}
                          </Stack>
                          <Typography variant="caption" sx={{ color: '#64748B', fontWeight: 600 }}>
                            ({meta.reviews})
                          </Typography>
                        </Stack>

                        {/* Fiyat & Rozet Satırı */}
                        <Box
                          sx={{
                            display: 'flex',
                            justifyContent: 'space-between',
                            alignItems: 'flex-end',
                            mt: 'auto',
                            mb: 2,
                          }}
                        >
                          <Box>
                            <Typography
                              variant="caption"
                              sx={{
                                color: '#94A3B8',
                                textDecoration: 'line-through',
                                display: 'block',
                                lineHeight: 1.2,
                                fontSize: '0.75rem',
                              }}
                            >
                              {simOldPrice.toLocaleString('tr-TR')} ₺
                            </Typography>
                            <Typography
                              variant="subtitle1"
                              sx={{
                                fontWeight: 800,
                                color: '#0F172A',
                                fontSize: '1.15rem',
                                lineHeight: 1.2,
                              }}
                            >
                              {simProd.price.toLocaleString('tr-TR')} ₺
                            </Typography>
                          </Box>

                          <Box
                            sx={{
                              bgcolor: '#EFF6FF',
                              color: '#2563EB',
                              fontWeight: 700,
                              fontSize: '0.72rem',
                              px: 1.2,
                              py: 0.4,
                              borderRadius: 1.2,
                            }}
                          >
                            {meta.tag}
                          </Box>
                        </Box>

                        {/* Sepete Ekle Butonu */}
                        <Button
                          variant="contained"
                          fullWidth
                          startIcon={simInCart ? <CheckIcon /> : <ShoppingCartOutlinedIcon />}
                          onClick={(e) => {
                            e.stopPropagation();
                            addtocart({
                              product_id: simProd.product_id,
                              product_name: simProd.product_name,
                              price: simProd.price,
                              stock: simProd.stock,
                              thumbnail_image: simProd.thumbnail_image,
                              quantity: 1,
                            });
                          }}
                          sx={{
                            bgcolor: '#0F275A',
                            color: '#FFFFFF',
                            py: 1.1,
                            fontWeight: 700,
                            fontSize: '0.875rem',
                            borderRadius: 2,
                            textTransform: 'none',
                            boxShadow: 'none',
                            '&:hover': {
                              bgcolor: '#1E3A8A',
                              boxShadow: '0 4px 12px rgba(15, 39, 90, 0.25)',
                            },
                          }}
                        >
                          {simInCart ? 'Sepete Eklendi' : 'Sepete Ekle'}
                        </Button>
                      </Box>
                    </Paper>
                  </Grid>
                );
              })}
            </Grid>
          </Box>
        )}

        {/* 4. GÜVEN & HİZMET BANNERI (Koyu Lacivert Zemin) */}
        <Box
          sx={{
            bgcolor: '#0F275A',
            borderRadius: 4,
            mt: 8,
            mb: 4,
            p: { xs: 4, md: 5 },
            color: '#FFFFFF',
          }}
        >
          <Grid container spacing={4} justifyContent="center">
            {/* 1. Ücretsiz Kargo */}
            <Grid size={{ xs: 12, sm: 6, md: 3 }} sx={{ textAlign: 'center' }}>
              <Box
                sx={{
                  display: 'inline-flex',
                  p: 2,
                  bgcolor: 'rgba(255, 255, 255, 0.12)',
                  borderRadius: '50%',
                  mb: 2,
                  color: '#FFFFFF',
                }}
              >
                <LocalShippingOutlinedIcon sx={{ fontSize: 36, color: '#FFFFFF' }} />
              </Box>
              <Typography
                variant="h6"
                gutterBottom
                sx={{
                  color: '#FFFFFF',
                  fontWeight: 800,
                  fontSize: '1.05rem',
                  fontFamily: '"Outfit", sans-serif',
                }}
              >
                Ücretsiz Kargo
              </Typography>
              <Typography
                variant="body2"
                sx={{
                  color: '#E2E8F0',
                  fontSize: '0.875rem',
                  lineHeight: 1.5,
                }}
              >
                Tüm siparişlerinizde ücretsiz kargo fırsatı.
              </Typography>
            </Grid>

            {/* 2. 256-bit SSL Güvenlik */}
            <Grid size={{ xs: 12, sm: 6, md: 3 }} sx={{ textAlign: 'center' }}>
              <Box
                sx={{
                  display: 'inline-flex',
                  p: 2,
                  bgcolor: 'rgba(255, 255, 255, 0.12)',
                  borderRadius: '50%',
                  mb: 2,
                  color: '#FFFFFF',
                }}
              >
                <SecurityOutlinedIcon sx={{ fontSize: 36, color: '#FFFFFF' }} />
              </Box>
              <Typography
                variant="h6"
                gutterBottom
                sx={{
                  color: '#FFFFFF',
                  fontWeight: 800,
                  fontSize: '1.05rem',
                  fontFamily: '"Outfit", sans-serif',
                }}
              >
                256-bit SSL Güvenlik
              </Typography>
              <Typography
                variant="body2"
                sx={{
                  color: '#E2E8F0',
                  fontSize: '0.875rem',
                  lineHeight: 1.5,
                }}
              >
                Uçtan uca şifreli ödeme altyapısıyla güvendesiniz.
              </Typography>
            </Grid>

            {/* 3. 14 Gün Koşulsuz İade */}
            <Grid size={{ xs: 12, sm: 6, md: 3 }} sx={{ textAlign: 'center' }}>
              <Box
                sx={{
                  display: 'inline-flex',
                  p: 2,
                  bgcolor: 'rgba(255, 255, 255, 0.12)',
                  borderRadius: '50%',
                  mb: 2,
                  color: '#FFFFFF',
                }}
              >
                <AssignmentReturnOutlinedIcon sx={{ fontSize: 36, color: '#FFFFFF' }} />
              </Box>
              <Typography
                variant="h6"
                gutterBottom
                sx={{
                  color: '#FFFFFF',
                  fontWeight: 800,
                  fontSize: '1.05rem',
                  fontFamily: '"Outfit", sans-serif',
                }}
              >
                14 Gün Koşulsuz İade
              </Typography>
              <Typography
                variant="body2"
                sx={{
                  color: '#E2E8F0',
                  fontSize: '0.875rem',
                  lineHeight: 1.5,
                }}
              >
                Memnun kalmazsanız anında iade ve değişim kolaylığı.
              </Typography>
            </Grid>

            {/* 4. 7/24 Canlı Destek */}
            <Grid size={{ xs: 12, sm: 6, md: 3 }} sx={{ textAlign: 'center' }}>
              <Box
                sx={{
                  display: 'inline-flex',
                  p: 2,
                  bgcolor: 'rgba(255, 255, 255, 0.12)',
                  borderRadius: '50%',
                  mb: 2,
                  color: '#FFFFFF',
                }}
              >
                <SupportAgentOutlinedIcon sx={{ fontSize: 36, color: '#FFFFFF' }} />
              </Box>
              <Typography
                variant="h6"
                gutterBottom
                sx={{
                  color: '#FFFFFF',
                  fontWeight: 800,
                  fontSize: '1.05rem',
                  fontFamily: '"Outfit", sans-serif',
                }}
              >
                7/24 Canlı Destek
              </Typography>
              <Typography
                variant="body2"
                sx={{
                  color: '#E2E8F0',
                  fontSize: '0.875rem',
                  lineHeight: 1.5,
                }}
              >
                Uzman ekibimiz her adımda size yardımcı olmaya hazır.
              </Typography>
            </Grid>
          </Grid>
        </Box>
      </Container>
    </Box>
  );
}

