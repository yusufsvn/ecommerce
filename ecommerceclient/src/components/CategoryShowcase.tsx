import { useState, useEffect } from 'react';
import { Box, Typography, Button, Container } from '@mui/material';
import { Grid } from '@mui/material';
import ArrowForwardIcon from '@mui/icons-material/ArrowForward';
import { useNavigate } from 'react-router';
import { categoryService } from '../Services/categoryService';
import type { GetCategories } from '../models/ProductResponse';

interface CategoryCardItem {
  id: string;
  title: string;
  subtitle: string;
  tag?: string;
  image: string;
  categoryQuery?: string;
}

const FEATURED_CATEGORIES: CategoryCardItem[] = [
  {
    id: 'furniture',
    title: 'Mobilya',
    subtitle: 'Furniture',
    tag: 'Yeni Sezon',
    image: 'https://images.unsplash.com/photo-1618221195710-dd6b41faaea6?auto=format&fit=crop&w=1000&q=80',
    categoryQuery: 'Mobilya'
  },
  {
    id: 'lighting',
    title: 'Aydınlatma',
    subtitle: 'Lighting',
    tag: 'Tasarım',
    image: 'https://images.unsplash.com/photo-1507473885765-e6ed057f782c?auto=format&fit=crop&w=1000&q=80',
    categoryQuery: 'Aydınlatma'
  },
  {
    id: 'accessories',
    title: 'Aksesuar',
    subtitle: 'Accessories',
    tag: 'Popüler',
    image: 'https://images.unsplash.com/photo-1612196808214-b8e1d6145a8c?auto=format&fit=crop&w=1000&q=80',
    categoryQuery: 'Aksesuar'
  },
  {
    id: 'kids',
    title: 'Çocuk',
    subtitle: 'Kids',
    tag: 'Doğal & Özel',
    image: 'https://images.unsplash.com/photo-1519689680058-324335c77eba?auto=format&fit=crop&w=1000&q=80',
    categoryQuery: 'Çocuk'
  }
];

export default function CategoryShowcase() {
  const navigate = useNavigate();
  const [dbCategories, setDbCategories] = useState<GetCategories[]>([]);

  useEffect(() => {
    const loadCategories = async () => {
      try {
        const data = await categoryService.getCategories();
        if (data && Array.isArray(data)) {
          setDbCategories(data);
        }
      } catch (err) {
        console.error('Kategoriler alınırken hata oluştu:', err);
      }
    };
    loadCategories();
  }, []);

  const handleCategoryClick = (item: CategoryCardItem) => {
    // Veritabanındaki kategori ile eşleşen bir ID var mı kontrol et
    const matched = dbCategories.find(
      (c) =>
        c.category_name.toLowerCase().trim() === item.title.toLowerCase().trim() ||
        c.category_name.toLowerCase().trim() === item.subtitle.toLowerCase().trim() ||
        (item.categoryQuery && c.category_name.toLowerCase().trim() === item.categoryQuery.toLowerCase().trim())
    );

    if (matched) {
      navigate(`/products?category=${encodeURIComponent(matched.category_id)}`);
    } else {
      // Eşleşme yoksa kategori adıyla filtrelemeyi dener veya ürünler sayfasına yönlendirir
      navigate(`/products?category=${encodeURIComponent(item.title)}`);
    }
  };

  return (
    <Container maxWidth="lg" sx={{ py: { xs: 5, md: 8 }, px: { xs: 2, md: 6 } }}>
      {/* Üst Başlık Bölümü */}
      <Box
        sx={{
          display: 'flex',
          flexDirection: { xs: 'column', sm: 'row' },
          justifyContent: 'space-between',
          alignItems: { xs: 'flex-start', sm: 'flex-end' },
          mb: { xs: 3, md: 4 },
          gap: 2,
        }}
      >
        <Box>
          <Typography
            variant="caption"
            sx={{
              color: 'secondary.main',
              fontWeight: 700,
              textTransform: 'uppercase',
              letterSpacing: '2px',
              fontSize: '0.85rem',
              display: 'inline-block',
              mb: 0.5,
            }}
          >
            Koleksiyonlar
          </Typography>
          <Typography
            variant="h4"
            component="h2"
            sx={{
              fontWeight: 800,
              color: 'primary.main',
              letterSpacing: '-0.5px',
            }}
          >
            Kategorileri Keşfedin
          </Typography>
          <Typography
            variant="body1"
            sx={{
              color: 'text.secondary',
              mt: 0.5,
              fontSize: { xs: '0.95rem', md: '1.05rem' },
            }}
          >
            Yaşam alanınıza ve tarzınıza ilham verecek seçkin tasarım koleksiyonları
          </Typography>
        </Box>

        <Button
          variant="text"
          color="primary"
          endIcon={<ArrowForwardIcon />}
          onClick={() => navigate('/products')}
          sx={{
            fontWeight: 700,
            fontSize: '0.95rem',
            px: 2,
            py: 1,
            borderRadius: 3,
            '&:hover': {
              bgcolor: 'rgba(30, 58, 138, 0.06)',
              transform: 'translateX(4px)',
            },
            transition: 'all 0.25s ease',
          }}
        >
          Tümünü Gör
        </Button>
      </Box>

      {/* 4 Dikey Büyük Resimli Kategori Kartı */}
      <Grid container spacing={{ xs: 2, sm: 2.5, md: 3 }}>
        {FEATURED_CATEGORIES.map((cat) => (
          <Grid size={{ xs: 12, sm: 6, md: 3 }} key={cat.id}>
            <Box
              onClick={() => handleCategoryClick(cat)}
              sx={{
                position: 'relative',
                height: { xs: 380, sm: 440, md: 500 },
                borderRadius: 4,
                overflow: 'hidden',
                cursor: 'pointer',
                border: '1px solid',
                borderColor: 'divider',
                bgcolor: 'grey.900',
                boxShadow: '0 4px 20px rgba(0, 0, 0, 0.07)',
                transition: 'all 0.4s cubic-bezier(0.16, 1, 0.3, 1)',
                '&:hover': {
                  transform: 'translateY(-8px)',
                  boxShadow: '0 20px 35px -10px rgba(30, 58, 138, 0.35)',
                  '& .category-img': {
                    transform: 'scale(1.08)',
                  },
                  '& .category-overlay': {
                    background:
                      'linear-gradient(180deg, rgba(0,0,0,0.05) 0%, rgba(0,0,0,0.2) 40%, rgba(15, 23, 42, 0.85) 75%, rgba(15, 23, 42, 0.96) 100%)',
                  },
                  '& .category-btn': {
                    transform: 'translateX(0)',
                    opacity: 1,
                    bgcolor: 'secondary.main',
                    color: '#FFFFFF',
                  },
                  '& .category-title': {
                    transform: 'translateY(-2px)',
                  },
                },
              }}
            >
              {/* Arka Plan Görseli */}
              <Box
                component="img"
                className="category-img"
                src={cat.image}
                alt={cat.title}
                loading="lazy"
                sx={{
                  position: 'absolute',
                  top: 0,
                  left: 0,
                  width: '100%',
                  height: '100%',
                  objectFit: 'cover',
                  objectPosition: 'center',
                  transition: 'transform 0.7s cubic-bezier(0.25, 1, 0.5, 1)',
                }}
              />

              {/* Degrade Karartma Katmanı (Alttaki yazıların net okunması için) */}
              <Box
                className="category-overlay"
                sx={{
                  position: 'absolute',
                  inset: 0,
                  background:
                    'linear-gradient(180deg, rgba(0,0,0,0.02) 0%, rgba(0,0,0,0.1) 45%, rgba(15, 23, 42, 0.75) 80%, rgba(15, 23, 42, 0.92) 100%)',
                  transition: 'background 0.4s ease',
                  zIndex: 1,
                }}
              />

              {/* Sağ Üst Etiket (Tag) */}
              {cat.tag && (
                <Box
                  sx={{
                    position: 'absolute',
                    top: 16,
                    right: 16,
                    zIndex: 2,
                    bgcolor: 'rgba(255, 255, 255, 0.88)',
                    backdropFilter: 'blur(8px)',
                    color: 'primary.dark',
                    px: 1.5,
                    py: 0.5,
                    borderRadius: 2,
                    fontSize: '0.75rem',
                    fontWeight: 700,
                    letterSpacing: '0.5px',
                    boxShadow: '0 2px 8px rgba(0,0,0,0.15)',
                  }}
                >
                  {cat.tag}
                </Box>
              )}

              {/* Alt Kısım: Kategori Başlığı ve Keşfet Butonu */}
              <Box
                sx={{
                  position: 'absolute',
                  bottom: 0,
                  left: 0,
                  right: 0,
                  p: { xs: 2.5, md: 3 },
                  zIndex: 2,
                  display: 'flex',
                  alignItems: 'flex-end',
                  justifyContent: 'space-between',
                }}
              >
                <Box sx={{ pr: 1 }}>
                  {/* İngilizce Alt Başlık / Vurgu */}
                  <Typography
                    variant="caption"
                    sx={{
                      color: 'rgba(255, 255, 255, 0.82)',
                      fontSize: '0.8rem',
                      fontWeight: 600,
                      textTransform: 'uppercase',
                      letterSpacing: '1.5px',
                      display: 'block',
                      mb: 0.3,
                    }}
                  >
                    {cat.subtitle}
                  </Typography>

                  {/* Ana Kategori Başlığı */}
                  <Typography
                    variant="h4"
                    className="category-title"
                    sx={{
                      fontWeight: 800,
                      color: '#FFFFFF',
                      fontSize: { xs: '1.65rem', sm: '1.85rem', md: '2.1rem' },
                      lineHeight: 1.1,
                      letterSpacing: '-0.5px',
                      transition: 'transform 0.3s ease',
                      textShadow: '0 2px 10px rgba(0, 0, 0, 0.4)',
                    }}
                  >
                    {cat.title}
                  </Typography>
                </Box>

                {/* Sağ Altta Şık Yuvarlak Keşfet Oku */}
                <Box
                  className="category-btn"
                  sx={{
                    width: 44,
                    height: 44,
                    borderRadius: '50%',
                    bgcolor: 'rgba(255, 255, 255, 0.25)',
                    backdropFilter: 'blur(10px)',
                    color: '#FFFFFF',
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'center',
                    flexShrink: 0,
                    boxShadow: '0 4px 12px rgba(0,0,0,0.2)',
                    transition: 'all 0.3s cubic-bezier(0.16, 1, 0.3, 1)',
                  }}
                >
                  <ArrowForwardIcon sx={{ fontSize: '1.25rem' }} />
                </Box>
              </Box>
            </Box>
          </Grid>
        ))}
      </Grid>
    </Container>
  );
}
