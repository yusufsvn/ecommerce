import React, { useState, useEffect, useMemo } from 'react';
import { useSearchParams, useNavigate } from 'react-router';
import {
  Container,
  Typography,
  Box,
  Paper,
  Button,
  IconButton,
  Grid,
  Slider,
  OutlinedInput,
  Select,
  MenuItem,
  Chip,
  Rating,
  Collapse,
  Breadcrumbs,
  Link as MuiLink,
  Tooltip,
  Skeleton,
} from '@mui/material';
import type { SelectChangeEvent } from '@mui/material';

// Icons
import TuneIcon from '@mui/icons-material/Tune';
import RefreshIcon from '@mui/icons-material/Refresh';
import ExpandLess from '@mui/icons-material/ExpandLess';
import ExpandMore from '@mui/icons-material/ExpandMore';
import GridViewIcon from '@mui/icons-material/GridView';
import ViewListIcon from '@mui/icons-material/ViewList';
import ShoppingCartOutlinedIcon from '@mui/icons-material/ShoppingCartOutlined';
import CheckIcon from '@mui/icons-material/Check';
import LocalShippingOutlinedIcon from '@mui/icons-material/LocalShippingOutlined';
import VerifiedUserOutlinedIcon from '@mui/icons-material/VerifiedUserOutlined';
import AutorenewOutlinedIcon from '@mui/icons-material/AutorenewOutlined';
import SupportAgentOutlinedIcon from '@mui/icons-material/SupportAgentOutlined';
import HomeOutlinedIcon from '@mui/icons-material/HomeOutlined';
import ArrowBackIosNewIcon from '@mui/icons-material/ArrowBackIosNew';
import ArrowForwardIosIcon from '@mui/icons-material/ArrowForwardIos';

import type { ProductListProps } from '../models/Props';
import { productService } from '../Services/productService';
import type { GetProducts } from '../models/ProductResponse';
import Spinner from './Spinner';

export interface CategoryData {
  category_id: string;
  category_name: string;
}

export default function ProductsPage(props: ProductListProps) {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const searchQuery = searchParams.get('searchQuery');
  const categoryParam = searchParams.get('category');

  const [products, setProducts] = useState<GetProducts[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [categories, setCategories] = useState<CategoryData[]>([]);
  const [categoriesLoading, setCategoriesLoading] = useState<boolean>(true);

  // Filtering states
  const [selectedCategory, setSelectedCategory] = useState<string | null>(null);
  const [appliedCategory, setAppliedCategory] = useState<string | null>(null);
  const [categoriesOpen, setCategoriesOpen] = useState<boolean>(true);

  // Price Range State (0 - 10000 TL)
  const [priceRange, setPriceRange] = useState<number[]>([0, 10000]);
  const [minPriceInput, setMinPriceInput] = useState<number>(0);
  const [maxPriceInput, setMaxPriceInput] = useState<number>(10000);
  const [appliedPriceRange, setAppliedPriceRange] = useState<number[]>([0, 10000]);

  // Toolbar states
  const [viewMode, setViewMode] = useState<'grid' | 'list'>('grid');
  const [sortMethod, setSortMethod] = useState<string>('newest');
  const [itemsPerPage, setItemsPerPage] = useState<number>(8);
  const [currentPage, setCurrentPage] = useState<number>(1);
  const [filterTrigger, setFilterTrigger] = useState<number>(0);

  // Fetch categories from backend API
  useEffect(() => {
    const fetchCats = async () => {
      setCategoriesLoading(true);
      try {
        const catData = await productService.getCategories();
        let rawList: any[] = [];
        if (Array.isArray(catData)) {
          rawList = catData;
        } else if (catData && Array.isArray((catData as any).data)) {
          rawList = (catData as any).data;
        } else if (catData && Array.isArray((catData as any).Data)) {
          rawList = (catData as any).Data;
        }

        const normalized: CategoryData[] = rawList
          .map((item: any) => ({
            category_id: item.category_id || item.CategoryId || item.id || '',
            category_name: item.category_name || item.CategoryName || item.name || '',
          }))
          .filter((c) => Boolean(c.category_name && c.category_name.trim()));

        setCategories(normalized);
      } catch (err) {
        console.error('Kategoriler alınırken hata oluştu:', err);
      } finally {
        setCategoriesLoading(false);
      }
    };
    fetchCats();
  }, []);

  // Sync category param from URL
  useEffect(() => {
    if (categoryParam) {
      setSelectedCategory(categoryParam);
      setAppliedCategory(categoryParam);
    }
  }, [categoryParam]);

  // Fetch products from backend whenever category, search, sort or filterTrigger changes
  useEffect(() => {
    const fetchProducts = async () => {
      setLoading(true);
      try {
        let datas;
        if (searchQuery) {
          datas = await productService.searchProduct(searchQuery, 100);
        } else if (appliedCategory) {
          // Kategori ID'sini belirle (ID veya isim ile eşleşebilir)
          const matchedCat = categories.find(
            (c) =>
              c.category_id.toLowerCase() === appliedCategory.toLowerCase() ||
              c.category_name.toLowerCase() === appliedCategory.toLowerCase()
          );
          const catId = matchedCat ? matchedCat.category_id : appliedCategory;
          datas = await productService.getFilteredProducts(100, 0, [catId], sortMethod || 'newest');
        } else {
          datas = await productService.getPublicProducts(0, 100);
        }

        let productList: GetProducts[] = [];
        if (datas) {
          if (Array.isArray(datas)) {
            productList = datas;
          } else if (Array.isArray((datas as any).Data)) {
            productList = (datas as any).Data;
          } else if (Array.isArray((datas as any).data)) {
            productList = (datas as any).data;
          }
        }
        setProducts(productList);
      } catch (err) {
        console.error('Ürünler yüklenirken hata oluştu:', err);
        setProducts([]);
      } finally {
        setLoading(false);
      }
    };
    fetchProducts();
  }, [searchQuery, appliedCategory, filterTrigger, sortMethod, categories]);

  // Filter and sort products
  const processedProducts = useMemo(() => {
    let list = [...products];

    // 1. Fiyat Filtresi
    const [minP, maxP] = appliedPriceRange;
    list = list.filter((p) => p.price >= minP && p.price <= maxP);

    // 2. Sıralama
    switch (sortMethod) {
      case 'newest':
        // En yeniler
        break;
      case 'priceAsc':
        list.sort((a, b) => a.price - b.price);
        break;
      case 'priceDesc':
        list.sort((a, b) => b.price - a.price);
        break;
      case 'popular':
        list.sort((a, b) => b.stock - a.stock);
        break;
      case 'nameAsc':
        list.sort((a, b) => a.product_name.localeCompare(b.product_name, 'tr'));
        break;
      case 'nameDesc':
        list.sort((a, b) => b.product_name.localeCompare(a.product_name, 'tr'));
        break;
      default:
        break;
    }

    return list;
  }, [products, appliedPriceRange, sortMethod]);

  // Pagination calculation
  const totalItems = processedProducts.length;
  const totalPages = Math.max(1, Math.ceil(totalItems / itemsPerPage));
  const startIndex = (currentPage - 1) * itemsPerPage;
  const currentProducts = processedProducts.slice(startIndex, startIndex + itemsPerPage);

  // Kategori adını id veya name üzerinden bulma
  const getCategoryDisplayName = (idOrName: string | null): string => {
    if (!idOrName) return '';
    const found = categories.find(
      (c) =>
        c.category_id.toLowerCase() === idOrName.toLowerCase() ||
        c.category_name.toLowerCase() === idOrName.toLowerCase()
    );
    return found ? found.category_name : idOrName;
  };

  // Ürünün kartında gösterilecek kategori adı
  const getProductCategoryTitle = (product: GetProducts): string => {
    const p = product as any;
    if (p.Categories && Array.isArray(p.Categories) && p.Categories.length > 0) {
      return p.Categories[0].category_name || p.Categories[0].name || '';
    }
    if (p.categories && Array.isArray(p.categories) && p.categories.length > 0) {
      return p.categories[0].category_name || p.categories[0].name || '';
    }
    if (p.category_name) {
      return p.category_name;
    }
    const pName = (product.product_name || '').toLowerCase();
    for (const cat of categories) {
      const cName = cat.category_name.toLowerCase();
      if (pName.includes(cName)) {
        return cat.category_name;
      }
    }
    if (appliedCategory) {
      return getCategoryDisplayName(appliedCategory);
    }
    return categories.length > 0 ? categories[0].category_name : 'Genel';
  };

  // Filter Actions
  const handleCategoryClick = (cat: CategoryData) => {
    const isSelected =
      selectedCategory?.toLowerCase() === cat.category_id.toLowerCase() ||
      selectedCategory?.toLowerCase() === cat.category_name.toLowerCase();

    if (isSelected) {
      setSelectedCategory(null);
      setAppliedCategory(null);
      navigate('/products', { replace: true });
    } else {
      setSelectedCategory(cat.category_id);
      setAppliedCategory(cat.category_id);
      navigate(`/products?category=${encodeURIComponent(cat.category_id)}`, { replace: true });
    }
    setCurrentPage(1);
    setFilterTrigger((prev) => prev + 1);
  };

  const handleSelectAllCategories = () => {
    setSelectedCategory(null);
    setAppliedCategory(null);
    navigate('/products', { replace: true });
    setCurrentPage(1);
    setFilterTrigger((prev) => prev + 1);
  };

  const handleApplyFilters = () => {
    setAppliedCategory(selectedCategory);
    setAppliedPriceRange([minPriceInput, maxPriceInput]);
    if (selectedCategory) {
      navigate(`/products?category=${encodeURIComponent(selectedCategory)}`, { replace: true });
    } else {
      navigate('/products', { replace: true });
    }
    setCurrentPage(1);
    setFilterTrigger((prev) => prev + 1);
  };

  const handleApplyPriceOnly = () => {
    setAppliedPriceRange([minPriceInput, maxPriceInput]);
    setCurrentPage(1);
  };

  const handleResetFilters = () => {
    setSelectedCategory(null);
    setAppliedCategory(null);
    setPriceRange([0, 10000]);
    setMinPriceInput(0);
    setMaxPriceInput(10000);
    setAppliedPriceRange([0, 10000]);
    navigate('/products', { replace: true });
    setCurrentPage(1);
    setFilterTrigger((prev) => prev + 1);
  };

  const handleSliderChange = (_event: Event, newValue: number | number[]) => {
    const arr = newValue as number[];
    setPriceRange(arr);
    setMinPriceInput(arr[0]);
    setMaxPriceInput(arr[1]);
  };

  const handleSortChange = (event: SelectChangeEvent) => {
    setSortMethod(event.target.value);
    setCurrentPage(1);
  };

  return (
    <Box sx={{ bgcolor: 'background.default', pb: 8 }}>
      <Container maxWidth="xl" sx={{ pt: 3, px: { xs: 2, md: 4 } }}>
        {/* ---------------- 1. ÜST BAŞLIK & BREADCRUMB ---------------- */}
        <Box sx={{ mb: 4 }}>
          <Breadcrumbs aria-label="breadcrumb" sx={{ mb: 1.5, fontSize: '0.85rem' }}>
            <MuiLink
              underline="hover"
              color="inherit"
              onClick={() => navigate('/')}
              sx={{ display: 'flex', alignItems: 'center', cursor: 'pointer' }}
            >
              <HomeOutlinedIcon sx={{ mr: 0.5, fontSize: '1.1rem' }} />
              Ana Sayfa
            </MuiLink>
            <Typography color={appliedCategory ? 'inherit' : 'text.primary'} sx={{ fontWeight: 600, fontSize: '0.85rem', cursor: appliedCategory ? 'pointer' : 'default' }} onClick={() => { if (appliedCategory) { handleSelectAllCategories(); } }}>
              Ürünler
            </Typography>
            {appliedCategory && (
              <Typography color="text.primary" sx={{ fontWeight: 700, fontSize: '0.85rem' }}>
                {getCategoryDisplayName(appliedCategory)}
              </Typography>
            )}
          </Breadcrumbs>

          <Box
            sx={{
              display: 'flex',
              flexDirection: { xs: 'column', sm: 'row' },
              alignItems: { xs: 'flex-start', sm: 'center' },
              justifyContent: 'space-between',
              gap: 2,
            }}
          >
            <Typography
              variant="h4"
              component="h1"
              sx={{
                fontFamily: '"Outfit", sans-serif',
                fontWeight: 800,
                color: 'primary.main',
                fontSize: { xs: '1.8rem', md: '2.3rem' },
                letterSpacing: '-0.5px',
              }}
            >
              {searchQuery
                ? `"${searchQuery}" için Sonuçlar`
                : appliedCategory
                  ? getCategoryDisplayName(appliedCategory)
                  : 'Tüm Ürünler'}
            </Typography>


          </Box>
        </Box>

        {/* ---------------- 2. ANA DÜZEN (SIDEBAR + İÇERİK) ---------------- */}
        <Grid container spacing={4}>
          {/* ================ SOL FİLTRELEME PANELİ ================ */}
          <Grid size={{ xs: 12, md: 3 }}>
            <Paper
              elevation={0}
              sx={{
                p: 3,
                borderRadius: 3,
                bgcolor: '#FFFFFF',
                border: '1px solid rgba(229, 231, 235, 0.8)',
                boxShadow: '0 1px 3px 0 rgba(0, 0, 0, 0.05), 0 1px 2px -1px rgba(0, 0, 0, 0.05)',
              }}
            >
              {/* Filtre Başlığı & Temizle Butonu */}
              <Box
                sx={{
                  display: 'flex',
                  alignItems: 'center',
                  justifyContent: 'space-between',
                  pb: 2,
                  mb: 2,
                  borderBottom: '1px solid #F3F4F6',
                }}
              >
                <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                  <TuneIcon sx={{ fontSize: 20, color: 'primary.main' }} />
                  <Typography
                    variant="subtitle1"
                    sx={{
                      fontFamily: '"Outfit", sans-serif',
                      fontWeight: 700,
                      color: 'primary.main',
                      fontSize: '1.1rem',
                    }}
                  >
                    Filtreleme
                  </Typography>
                </Box>

                <Button
                  size="small"
                  startIcon={<RefreshIcon sx={{ fontSize: 15 }} />}
                  onClick={handleResetFilters}
                  sx={{
                    color: 'secondary.main',
                    fontWeight: 700,
                    fontSize: '0.8rem',
                    p: 0.5,
                    '&:hover': { bgcolor: 'transparent', color: 'secondary.dark' },
                  }}
                >
                  Temizle
                </Button>
              </Box>

              {/* 1. Kategoriler Bölümü */}
              <Box sx={{ mb: 3 }}>
                <Box
                  onClick={() => setCategoriesOpen(!categoriesOpen)}
                  sx={{
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'space-between',
                    cursor: 'pointer',
                    py: 1,
                    mb: 1,
                  }}
                >
                  <Typography
                    variant="subtitle2"
                    sx={{
                      fontFamily: '"Outfit", sans-serif',
                      fontWeight: 700,
                      color: 'primary.main',
                      fontSize: '0.95rem',
                    }}
                  >
                    Kategoriler
                  </Typography>
                  {categoriesOpen ? (
                    <ExpandLess sx={{ fontSize: 20, color: 'text.secondary' }} />
                  ) : (
                    <ExpandMore sx={{ fontSize: 20, color: 'text.secondary' }} />
                  )}
                </Box>

                <Collapse in={categoriesOpen}>
                  <Box sx={{ display: 'flex', flexDirection: 'column', gap: 0.75, pl: 0.5 }}>
                    {/* Tüm Kategoriler Seçeneği */}
                    <Box
                      onClick={handleSelectAllCategories}
                      sx={{
                        px: 1.5,
                        py: 1,
                        borderRadius: 2,
                        cursor: 'pointer',
                        bgcolor: !selectedCategory ? '#EFF6FF' : 'transparent',
                        color: !selectedCategory ? 'primary.main' : 'text.primary',
                        fontWeight: !selectedCategory ? 700 : 500,
                        fontSize: '0.9rem',
                        display: 'flex',
                        alignItems: 'center',
                        justifyContent: 'space-between',
                        transition: 'all 0.2s ease',
                        border: !selectedCategory ? '1px solid #BFDBFE' : '1px solid transparent',
                        '&:hover': {
                          bgcolor: !selectedCategory ? '#EFF6FF' : '#F9FAFB',
                          color: 'primary.main',
                        },
                      }}
                    >
                      <Typography variant="body2" sx={{ fontWeight: 'inherit', color: 'inherit' }}>
                        Tüm Kategoriler
                      </Typography>
                    </Box>

                    {/* API'den Gelen Kategoriler */}
                    {categoriesLoading ? (
                      Array.from({ length: 5 }).map((_, i) => (
                        <Skeleton
                          key={i}
                          variant="rounded"
                          height={36}
                          sx={{ borderRadius: 2, bgcolor: '#F3F4F6' }}
                        />
                      ))
                    ) : categories.length > 0 ? (
                      categories.map((cat) => {
                        const isSelected =
                          selectedCategory?.toLowerCase() === cat.category_id.toLowerCase() ||
                          selectedCategory?.toLowerCase() === cat.category_name.toLowerCase();

                        return (
                          <Box
                            key={cat.category_id || cat.category_name}
                            onClick={() => handleCategoryClick(cat)}
                            sx={{
                              px: 1.5,
                              py: 1,
                              borderRadius: 2,
                              cursor: 'pointer',
                              bgcolor: isSelected ? '#EFF6FF' : 'transparent',
                              color: isSelected ? 'primary.main' : 'text.primary',
                              fontWeight: isSelected ? 700 : 500,
                              fontSize: '0.9rem',
                              display: 'flex',
                              alignItems: 'center',
                              justifyContent: 'space-between',
                              transition: 'all 0.2s ease',
                              border: isSelected ? '1px solid #BFDBFE' : '1px solid transparent',
                              '&:hover': {
                                bgcolor: isSelected ? '#EFF6FF' : '#F9FAFB',
                                color: 'primary.main',
                              },
                            }}
                          >
                            <Typography variant="body2" sx={{ fontWeight: 'inherit', color: 'inherit' }}>
                              {cat.category_name}
                            </Typography>
                          </Box>
                        );
                      })
                    ) : (
                      <Typography variant="caption" sx={{ color: 'text.secondary', px: 1.5, py: 1 }}>
                        Kategori bulunamadı
                      </Typography>
                    )}
                  </Box>
                </Collapse>
              </Box>

              {/* 2. Fiyat Aralığı Bölümü */}
              <Box sx={{ pt: 2, borderTop: '1px solid #F3F4F6', mb: 3 }}>
                <Typography
                  variant="subtitle2"
                  sx={{
                    fontFamily: '"Outfit", sans-serif',
                    fontWeight: 700,
                    color: 'primary.main',
                    fontSize: '0.95rem',
                    mb: 1.5,
                  }}
                >
                  Fiyat Aralığı
                </Typography>

                {/* Slider */}
                <Slider
                  value={priceRange}
                  onChange={handleSliderChange}
                  valueLabelDisplay="auto"
                  min={0}
                  max={10000}
                  step={100}
                  sx={{
                    color: 'primary.main',
                    mb: 2,
                    '& .MuiSlider-thumb': {
                      width: 18,
                      height: 18,
                      bgcolor: '#FFFFFF',
                      border: '3px solid #1E3A8A',
                      '&:hover, &.Mui-focusVisible': {
                        boxShadow: '0 0 0 8px rgba(30, 58, 138, 0.16)',
                      },
                    },
                    '& .MuiSlider-track': { bgcolor: 'primary.main' },
                    '& .MuiSlider-rail': { bgcolor: '#E5E7EB' },
                  }}
                />

                {/* Min / Max Etiketleri */}
                <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 1.5 }}>
                  <Typography variant="caption" sx={{ color: 'text.secondary', fontWeight: 600 }}>
                    0 ₺
                  </Typography>
                  <Typography variant="caption" sx={{ color: 'text.secondary', fontWeight: 600 }}>
                    10.000+ ₺
                  </Typography>
                </Box>

                {/* Min / Max Sayı Girdileri */}
                <Box sx={{ display: 'flex', gap: 1.5, mb: 2 }}>
                  <Box sx={{ flex: 1 }}>
                    <Typography variant="caption" sx={{ color: 'text.secondary', fontWeight: 600, display: 'block', mb: 0.5 }}>
                      Min (₺)
                    </Typography>
                    <OutlinedInput
                      size="small"
                      type="number"
                      value={minPriceInput}
                      onChange={(e) => {
                        const v = Number(e.target.value);
                        setMinPriceInput(v);
                        setPriceRange([v, maxPriceInput]);
                      }}
                      sx={{
                        borderRadius: '8px',
                        fontSize: '0.85rem',
                        '& fieldset': { borderColor: '#E5E7EB' },
                      }}
                    />
                  </Box>

                  <Box sx={{ flex: 1 }}>
                    <Typography variant="caption" sx={{ color: 'text.secondary', fontWeight: 600, display: 'block', mb: 0.5 }}>
                      Max (₺)
                    </Typography>
                    <OutlinedInput
                      size="small"
                      type="number"
                      value={maxPriceInput}
                      onChange={(e) => {
                        const v = Number(e.target.value);
                        setMaxPriceInput(v);
                        setPriceRange([minPriceInput, v]);
                      }}
                      sx={{
                        borderRadius: '8px',
                        fontSize: '0.85rem',
                        '& fieldset': { borderColor: '#E5E7EB' },
                      }}
                    />
                  </Box>
                </Box>

                {/* Fiyatı Uygula Butonu */}
                <Button
                  fullWidth
                  variant="outlined"
                  onClick={handleApplyPriceOnly}
                  sx={{
                    bgcolor: '#EFF6FF',
                    borderColor: '#BFDBFE',
                    color: 'primary.main',
                    fontWeight: 700,
                    fontSize: '0.85rem',
                    py: 0.8,
                    borderRadius: '8px',
                    '&:hover': {
                      bgcolor: '#DBEAFE',
                      borderColor: '#93C5FD',
                    },
                  }}
                >
                  Fiyatı Uygula
                </Button>
              </Box>

              {/* Filtreleri Uygula Ana Buton */}
              <Button
                fullWidth
                variant="contained"
                startIcon={<TuneIcon />}
                onClick={handleApplyFilters}
                sx={{
                  bgcolor: 'primary.main',
                  color: '#FFFFFF',
                  fontWeight: 700,
                  fontSize: '0.95rem',
                  py: 1.3,
                  borderRadius: '10px',
                  boxShadow: '0 4px 12px rgba(30, 58, 138, 0.25)',
                  '&:hover': {
                    bgcolor: 'primary.dark',
                    boxShadow: '0 6px 16px rgba(30, 58, 138, 0.35)',
                  },
                }}
              >
                Filtreleri Uygula
              </Button>
            </Paper>
          </Grid>

          {/* ================ SAĞ ÜRÜNLER & ARAÇ ÇUBUĞU ================ */}
          <Grid size={{ xs: 12, md: 9 }}>
            {/* 1. ÜST ARAÇ ÇUBUĞU (TOOLBAR) */}
            <Paper
              elevation={0}
              sx={{
                p: 2,
                mb: 3,
                borderRadius: 3,
                bgcolor: '#FFFFFF',
                border: '1px solid rgba(229, 231, 235, 0.8)',
                boxShadow: '0 1px 3px 0 rgba(0, 0, 0, 0.05)',
              }}
            >
              <Box
                sx={{
                  display: 'flex',
                  flexWrap: 'wrap',
                  alignItems: 'center',
                  justifyContent: 'space-between',
                  gap: 2,
                }}
              >
                {/* Sol: Görünüm Butonları & Sayaç */}
                <Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
                  <Box
                    sx={{
                      display: 'flex',
                      bgcolor: '#F3F4F6',
                      borderRadius: 2,
                      p: 0.3,
                      border: '1px solid #E5E7EB',
                    }}
                  >
                    <IconButton
                      size="small"
                      onClick={() => setViewMode('grid')}
                      sx={{
                        borderRadius: 1.5,
                        bgcolor: viewMode === 'grid' ? '#FFFFFF' : 'transparent',
                        color: viewMode === 'grid' ? 'primary.main' : 'text.secondary',
                        boxShadow: viewMode === 'grid' ? '0 1px 3px rgba(0,0,0,0.1)' : 'none',
                        p: 0.6,
                      }}
                    >
                      <GridViewIcon fontSize="small" />
                    </IconButton>
                    <IconButton
                      size="small"
                      onClick={() => setViewMode('list')}
                      sx={{
                        borderRadius: 1.5,
                        bgcolor: viewMode === 'list' ? '#FFFFFF' : 'transparent',
                        color: viewMode === 'list' ? 'primary.main' : 'text.secondary',
                        boxShadow: viewMode === 'list' ? '0 1px 3px rgba(0,0,0,0.1)' : 'none',
                        p: 0.6,
                      }}
                    >
                      <ViewListIcon fontSize="small" />
                    </IconButton>
                  </Box>

                  <Typography variant="body2" sx={{ color: 'text.secondary', fontSize: '0.875rem' }}>
                    <strong style={{ color: '#1E3A8A' }}>
                      {totalItems === 0 ? 0 : startIndex + 1} - {Math.min(startIndex + itemsPerPage, totalItems)}
                    </strong>{' '}
                    arası gösteriliyor (Toplam {totalItems})
                  </Typography>
                </Box>

                {/* Sağ: Sıralama & Göster Seçimi */}
                <Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
                  {/* Sıralama */}
                  <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                    <Typography variant="caption" sx={{ color: 'text.secondary', fontWeight: 600 }}>
                      Sıralama:
                    </Typography>
                    <Select
                      size="small"
                      value={sortMethod}
                      onChange={handleSortChange}
                      sx={{
                        borderRadius: '8px',
                        fontSize: '0.85rem',
                        fontWeight: 600,
                        bgcolor: '#FFFFFF',
                        minWidth: 150,
                        '& .MuiOutlinedInput-notchedOutline': { borderColor: '#E5E7EB' },
                      }}
                    >
                      <MenuItem value="newest">En Yeniler</MenuItem>
                      <MenuItem value="priceAsc">Fiyat (Düşükten Yükseğe)</MenuItem>
                      <MenuItem value="priceDesc">Fiyat (Yüksekten Düşüğe)</MenuItem>
                      <MenuItem value="popular">En Çok Satanlar</MenuItem>
                      <MenuItem value="nameAsc">A - Z</MenuItem>
                      <MenuItem value="nameDesc">Z - A</MenuItem>
                    </Select>
                  </Box>

                  {/* Sayfa Başına Göster */}
                  <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                    <Typography variant="caption" sx={{ color: 'text.secondary', fontWeight: 600 }}>
                      Göster:
                    </Typography>
                    <Box sx={{ display: 'flex', gap: 0.5 }}>
                      {[8, 16, 32].map((num) => (
                        <Button
                          key={num}
                          size="small"
                          onClick={() => {
                            setItemsPerPage(num);
                            setCurrentPage(1);
                          }}
                          sx={{
                            minWidth: 28,
                            height: 28,
                            p: 0,
                            borderRadius: '6px',
                            fontSize: '0.78rem',
                            fontWeight: 700,
                            bgcolor: itemsPerPage === num ? '#EFF6FF' : 'transparent',
                            color: itemsPerPage === num ? 'primary.main' : 'text.secondary',
                            border: itemsPerPage === num ? '1px solid #BFDBFE' : '1px solid transparent',
                            '&:hover': { bgcolor: '#F3F4F6' },
                          }}
                        >
                          {num}
                        </Button>
                      ))}
                    </Box>
                  </Box>
                </Box>
              </Box>

              {/* Aktif Filtre Rozetleri (Chips) */}
              {(appliedCategory || appliedPriceRange[0] > 0 || appliedPriceRange[1] < 10000) && (
                <Box
                  sx={{
                    display: 'flex',
                    flexWrap: 'wrap',
                    alignItems: 'center',
                    gap: 1,
                    mt: 2,
                    pt: 1.5,
                    borderTop: '1px solid #F3F4F6',
                  }}
                >
                  <Typography variant="caption" sx={{ fontWeight: 800, color: 'text.secondary', mr: 0.5 }}>
                    AKTİF:
                  </Typography>

                  {appliedCategory && (
                    <Chip
                      label={`Kategori: ${getCategoryDisplayName(appliedCategory)}`}
                      onDelete={() => {
                        setSelectedCategory(null);
                        setAppliedCategory(null);
                        navigate('/products', { replace: true });
                        setFilterTrigger((prev) => prev + 1);
                      }}
                      size="small"
                      sx={{
                        bgcolor: '#EFF6FF',
                        border: '1px solid #BFDBFE',
                        color: 'primary.main',
                        fontWeight: 600,
                        fontSize: '0.78rem',
                        '& .MuiChip-deleteIcon': { color: 'primary.main' },
                      }}
                    />
                  )}

                  {(appliedPriceRange[0] > 0 || appliedPriceRange[1] < 10000) && (
                    <Chip
                      label={`Fiyat: ${appliedPriceRange[0].toLocaleString('tr-TR')} ₺ - ${appliedPriceRange[1].toLocaleString('tr-TR')} ₺`}
                      onDelete={() => {
                        setPriceRange([0, 10000]);
                        setMinPriceInput(0);
                        setMaxPriceInput(10000);
                        setAppliedPriceRange([0, 10000]);
                      }}
                      size="small"
                      sx={{
                        bgcolor: '#EFF6FF',
                        border: '1px solid #BFDBFE',
                        color: 'primary.main',
                        fontWeight: 600,
                        fontSize: '0.78rem',
                        '& .MuiChip-deleteIcon': { color: 'primary.main' },
                      }}
                    />
                  )}

                  <Typography
                    variant="caption"
                    onClick={handleResetFilters}
                    sx={{
                      color: 'secondary.main',
                      fontWeight: 700,
                      cursor: 'pointer',
                      ml: 1,
                      '&:hover': { textDecoration: 'underline' },
                    }}
                  >
                    Tümünü Kaldır
                  </Typography>
                </Box>
              )}
            </Paper>

            {/* 2. ÜRÜN IZGARASI (GRID YA DA LİSTE) */}
            {loading ? (
              <Box sx={{ py: 12, display: 'flex', justifyContent: 'center' }}>
                <Spinner message="Ürünler listeleniyor..." />
              </Box>
            ) : currentProducts.length > 0 ? (
              <Grid container spacing={3}>
                {currentProducts.map((product, index) => {
                  const inCart = props.cartItems
                    ? props.cartItems.some((item) => item.product_id === product.product_id)
                    : false;

                  // Dinamik rozetler
                  const badges = [
                    { text: '%15 İndirim', color: '#F97316' },
                    { text: '%20 İndirim', color: '#F97316' },
                    { text: 'Yeni Sezon', color: '#1E3A8A' },
                    { text: 'Fırsat', color: '#EA580C' },
                    { text: '%25 İndirim', color: '#F97316' },
                    { text: 'Sınırlı Stok', color: '#172554' },
                  ];
                  const subBadges = ['Çok Satan', '', '', '', 'Kargo Bedava', ''];

                  const mainBadge = badges[index % badges.length];
                  const subBadge = subBadges[index % subBadges.length];

                  return (
                    <Grid
                      size={viewMode === 'grid' ? { xs: 12, sm: 6, md: 4 } : { xs: 12 }}
                      key={product.product_id}
                    >
                      <Paper
                        elevation={0}
                        onClick={() => navigate(`/product/${product.product_id}`)}
                        sx={{
                          height: '100%',
                          display: 'flex',
                          flexDirection: viewMode === 'grid' ? 'column' : { xs: 'column', sm: 'row' },
                          borderRadius: 3,
                          bgcolor: '#FFFFFF',
                          border: '1px solid rgba(229, 231, 235, 0.8)',
                          boxShadow: '0 1px 3px 0 rgba(0, 0, 0, 0.05)',
                          overflow: 'hidden',
                          cursor: 'pointer',
                          transition: 'all 0.3s cubic-bezier(0.16, 1, 0.3, 1)',
                          '&:hover': {
                            transform: 'translateY(-6px)',
                            boxShadow: '0 12px 24px -4px rgba(30, 58, 138, 0.12)',
                            borderColor: '#BFDBFE',
                            '& .product-img': {
                              transform: 'scale(1.05)',
                            },
                          },
                        }}
                      >
                        {/* Görsel Alanı */}
                        <Box
                          sx={{
                            position: 'relative',
                            height: viewMode === 'grid' ? 260 : 200,
                            width: viewMode === 'grid' ? '100%' : { sm: 220 },
                            flexShrink: 0,
                            bgcolor: '#FFFFFF',
                            display: 'flex',
                            alignItems: 'center',
                            justifyContent: 'center',
                            p: 2.5,
                            borderBottom: viewMode === 'grid' ? '1px solid #F3F4F6' : 'none',
                            borderRight: viewMode === 'grid' ? 'none' : { sm: '1px solid #F3F4F6' },
                          }}
                        >
                          {/* Sol Üst Rozetler */}
                          <Box sx={{ position: 'absolute', top: 12, left: 12, zIndex: 2, display: 'flex', flexDirection: 'column', gap: 0.5 }}>
                            {mainBadge && (
                              <Box
                                sx={{
                                  bgcolor: mainBadge.color,
                                  color: '#FFFFFF',
                                  px: 1.2,
                                  py: 0.4,
                                  borderRadius: 1.5,
                                  fontSize: '0.72rem',
                                  fontWeight: 800,
                                  letterSpacing: '0.3px',
                                  boxShadow: '0 2px 6px rgba(0,0,0,0.15)',
                                }}
                              >
                                {mainBadge.text}
                              </Box>
                            )}

                            {subBadge && (
                              <Box
                                sx={{
                                  bgcolor: '#1E3A8A',
                                  color: '#FFFFFF',
                                  px: 1,
                                  py: 0.3,
                                  borderRadius: 1.5,
                                  fontSize: '0.68rem',
                                  fontWeight: 700,
                                }}
                              >
                                {subBadge}
                              </Box>
                            )}
                          </Box>

                          {/* Fotoğraf */}
                          <Box
                            component="img"
                            className="product-img"
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

                        {/* Bilgi & Fiyat Alanı */}
                        <Box
                          sx={{
                            p: 2.5,
                            display: 'flex',
                            flexDirection: 'column',
                            flexGrow: 1,
                            justifyContent: 'space-between',
                          }}
                        >
                          <Box>
                            {/* Kategori & Yıldız Satırı */}
                            <Box
                              sx={{
                                display: 'flex',
                                alignItems: 'center',
                                justifyContent: 'space-between',
                                mb: 1,
                              }}
                            >
                              <Typography
                                variant="caption"
                                sx={{
                                  fontFamily: '"Inter", sans-serif',
                                  fontWeight: 600,
                                  color: '#6B7280',
                                  fontSize: '0.75rem',
                                }}
                              >
                                {getProductCategoryTitle(product)}
                              </Typography>

                              <Box sx={{ display: 'flex', alignItems: 'center', gap: 0.5 }}>
                                <Rating
                                  value={4.9}
                                  readOnly
                                  precision={0.1}
                                  size="small"
                                  sx={{ fontSize: '0.9rem', color: '#F59E0B' }}
                                />
                                <Typography variant="caption" sx={{ fontWeight: 700, color: '#374151', fontSize: '0.78rem' }}>
                                  4.9
                                </Typography>
                                <Typography variant="caption" sx={{ color: 'text.disabled', fontSize: '0.72rem' }}>
                                  ({(index + 1) * 35})
                                </Typography>
                              </Box>
                            </Box>

                            {/* Ürün Başlığı */}
                            <Typography
                              variant="subtitle1"
                              sx={{
                                fontFamily: '"Outfit", sans-serif',
                                fontWeight: 700,
                                color: 'primary.main',
                                fontSize: '1rem',
                                lineHeight: 1.35,
                                mb: 2,
                                display: '-webkit-box',
                                WebkitLineClamp: 2,
                                WebkitBoxOrient: 'vertical',
                                overflow: 'hidden',
                                height: 42,
                              }}
                            >
                              {product.product_name}
                            </Typography>
                          </Box>

                          {/* Fiyat ve Sepet Butonu */}
                          <Box
                            sx={{
                              display: 'flex',
                              alignItems: 'flex-end',
                              justifyContent: 'space-between',
                              pt: 1.5,
                              borderTop: '1px solid #F3F4F6',
                            }}
                          >
                            <Box>
                              <Typography
                                variant="caption"
                                sx={{
                                  color: 'text.disabled',
                                  textDecoration: 'line-through',
                                  fontSize: '0.78rem',
                                  display: 'block',
                                  lineHeight: 1,
                                  mb: 0.4,
                                }}
                              >
                                {(product.price * 1.25).toLocaleString('tr-TR', { minimumFractionDigits: 2 })} ₺
                              </Typography>
                              <Typography
                                variant="h6"
                                sx={{
                                  fontFamily: '"Outfit", sans-serif',
                                  fontWeight: 800,
                                  color: 'primary.main',
                                  fontSize: '1.25rem',
                                  lineHeight: 1,
                                }}
                              >
                                {product.price.toLocaleString('tr-TR', { minimumFractionDigits: 2 })} ₺
                              </Typography>
                            </Box>

                            {/* Sepet Butonu (Koyu Mavi Kare) */}
                            <Tooltip title={inCart ? 'Sepette Var' : 'Sepete Ekle'}>
                              <IconButton
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
                                  width: 42,
                                  height: 42,
                                  borderRadius: '10px',
                                  bgcolor: inCart ? 'success.main' : 'primary.main',
                                  color: '#FFFFFF',
                                  boxShadow: '0 4px 10px rgba(30, 58, 138, 0.25)',
                                  '&:hover': {
                                    bgcolor: inCart ? 'success.dark' : 'primary.dark',
                                    transform: 'scale(1.04)',
                                  },
                                  '&.Mui-disabled': {
                                    bgcolor: 'success.main',
                                    color: '#FFFFFF',
                                    opacity: 0.9,
                                  },
                                  transition: 'all 0.2s ease',
                                }}
                              >
                                {inCart ? <CheckIcon sx={{ fontSize: 20 }} /> : <ShoppingCartOutlinedIcon sx={{ fontSize: 20 }} />}
                              </IconButton>
                            </Tooltip>
                          </Box>
                        </Box>
                      </Paper>
                    </Grid>
                  );
                })}
              </Grid>
            ) : (
              <Paper
                elevation={0}
                sx={{
                  py: 10,
                  px: 4,
                  textAlign: 'center',
                  borderRadius: 3,
                  bgcolor: '#FFFFFF',
                  border: '1px solid rgba(229, 231, 235, 0.8)',
                }}
              >
                <Typography variant="h6" sx={{ color: 'primary.main', fontWeight: 700, mb: 1 }}>
                  Aradığınız kriterlere uygun ürün bulunamadı.
                </Typography>
                <Typography variant="body2" sx={{ color: 'text.secondary', mb: 3 }}>
                  Lütfen farklı bir filtre veya arama terimi deneyin.
                </Typography>
                <Button variant="contained" onClick={handleResetFilters} sx={{ bgcolor: 'primary.main' }}>
                  Filtreleri Sıfırla
                </Button>
              </Paper>
            )}

            {/* 3. SAYFALAMA (PAGINATION) */}
            {totalPages > 1 && (
              <Box
                sx={{
                  display: 'flex',
                  alignItems: 'center',
                  justifyContent: 'space-between',
                  flexWrap: 'wrap',
                  gap: 2,
                  mt: 5,
                  p: 2,
                  borderRadius: 3,
                  bgcolor: '#FFFFFF',
                  border: '1px solid rgba(229, 231, 235, 0.8)',
                }}
              >
                <Typography variant="body2" sx={{ color: 'text.secondary', fontWeight: 600 }}>
                  Sayfa <strong style={{ color: '#1E3A8A' }}>{currentPage}</strong> / {totalPages}
                </Typography>

                <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                  {/* Önceki Sayfa */}
                  <IconButton
                    size="small"
                    disabled={currentPage === 1}
                    onClick={() => {
                      setCurrentPage((p) => Math.max(1, p - 1));
                      window.scrollTo({ top: 0, behavior: 'smooth' });
                    }}
                    sx={{
                      borderRadius: 1.5,
                      border: '1px solid #E5E7EB',
                      bgcolor: '#FFFFFF',
                      width: 36,
                      height: 36,
                    }}
                  >
                    <ArrowBackIosNewIcon sx={{ fontSize: 13 }} />
                  </IconButton>

                  {/* Sayfa Numaraları */}
                  {Array.from({ length: totalPages }, (_, i) => i + 1).map((pageNum) => (
                    <Button
                      key={pageNum}
                      size="small"
                      onClick={() => {
                        setCurrentPage(pageNum);
                        window.scrollTo({ top: 0, behavior: 'smooth' });
                      }}
                      sx={{
                        minWidth: 36,
                        height: 36,
                        p: 0,
                        borderRadius: '8px',
                        fontSize: '0.85rem',
                        fontWeight: 700,
                        bgcolor: currentPage === pageNum ? 'primary.main' : 'transparent',
                        color: currentPage === pageNum ? '#FFFFFF' : 'text.primary',
                        border: currentPage === pageNum ? 'none' : '1px solid #E5E7EB',
                        '&:hover': {
                          bgcolor: currentPage === pageNum ? 'primary.dark' : '#F3F4F6',
                        },
                      }}
                    >
                      {pageNum}
                    </Button>
                  ))}

                  {/* Sonraki Sayfa */}
                  <IconButton
                    size="small"
                    disabled={currentPage === totalPages}
                    onClick={() => {
                      setCurrentPage((p) => Math.min(totalPages, p + 1));
                      window.scrollTo({ top: 0, behavior: 'smooth' });
                    }}
                    sx={{
                      borderRadius: 1.5,
                      border: '1px solid #E5E7EB',
                      bgcolor: '#FFFFFF',
                      width: 36,
                      height: 36,
                    }}
                  >
                    <ArrowForwardIosIcon sx={{ fontSize: 13 }} />
                  </IconButton>
                </Box>
              </Box>
            )}
          </Grid>
        </Grid>

        {/* ---------------- 3. ALT GÜVEN & DEĞER ROZETLERİ (4'LÜ KART) ---------------- */}
        <Box sx={{ mt: 8 }}>
          <Paper
            elevation={0}
            sx={{
              p: { xs: 3, md: 4 },
              borderRadius: 3,
              bgcolor: '#FFFFFF',
              border: '1px solid rgba(229, 231, 235, 0.8)',
              boxShadow: '0 1px 3px 0 rgba(0, 0, 0, 0.05)',
            }}
          >
            <Grid container spacing={3}>
              {[
                {
                  icon: <LocalShippingOutlinedIcon sx={{ fontSize: 26, color: 'primary.main' }} />,
                  title: 'Ücretsiz Hızlı Teslimat',
                  desc: '500 TL ve üzeri tüm siparişlerde',
                },
                {
                  icon: <VerifiedUserOutlinedIcon sx={{ fontSize: 26, color: 'primary.main' }} />,
                  title: 'Güvenli Ödeme Altyapısı',
                  desc: '256-Bit SSL şifreleme ve 3D Secure',
                },
                {
                  icon: <AutorenewOutlinedIcon sx={{ fontSize: 26, color: 'primary.main' }} />,
                  title: '14 Gün Kolay İade',
                  desc: 'Koşulsuz iade ve değişim hakkı',
                },
                {
                  icon: <SupportAgentOutlinedIcon sx={{ fontSize: 26, color: 'primary.main' }} />,
                  title: '7/24 Müşteri Desteği',
                  desc: 'Her an yanınızda canlı yardım',
                },
              ].map((item, idx) => (
                <Grid size={{ xs: 12, sm: 6, md: 3 }} key={idx}>
                  <Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
                    <Box
                      sx={{
                        width: 48,
                        height: 48,
                        borderRadius: 2.5,
                        bgcolor: '#EFF6FF',
                        display: 'flex',
                        alignItems: 'center',
                        justifyContent: 'center',
                        flexShrink: 0,
                        border: '1px solid #BFDBFE',
                      }}
                    >
                      {item.icon}
                    </Box>
                    <Box>
                      <Typography
                        variant="subtitle2"
                        sx={{
                          fontFamily: '"Outfit", sans-serif',
                          fontWeight: 700,
                          color: 'primary.main',
                          fontSize: '0.95rem',
                          lineHeight: 1.2,
                        }}
                      >
                        {item.title}
                      </Typography>
                      <Typography
                        variant="caption"
                        sx={{
                          fontFamily: '"Inter", sans-serif',
                          color: 'text.secondary',
                          fontSize: '0.78rem',
                          display: 'block',
                          mt: 0.3,
                        }}
                      >
                        {item.desc}
                      </Typography>
                    </Box>
                  </Box>
                </Grid>
              ))}
            </Grid>
          </Paper>
        </Box>
      </Container>
    </Box>
  );
}
