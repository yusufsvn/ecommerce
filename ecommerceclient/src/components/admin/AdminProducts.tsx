import { useEffect, useState } from 'react';
import { Box, Paper, Typography, Button, Stack, Container, Avatar, Chip, Pagination, TextField, InputAdornment, IconButton } from '@mui/material';
import EditIcon from '@mui/icons-material/Edit';
import DeleteIcon from '@mui/icons-material/Delete';
import SearchIcon from '@mui/icons-material/Search';
import { productService } from '../../Services/productService';
import type { GetCategoryWithProduct } from '../../models/ProductResponse';
import QuickEdit from './Pop-ups/quickedit';
import { toast } from 'react-toastify';
import Spinner from '../Spinner';
import { useNavigate } from 'react-router';

function AdminProducts() {
    const [products, setProducts] = useState<GetCategoryWithProduct[]>([]);
    const [loading, setLoading] = useState(true);
    const [selectedProductId, setSelectedProductId] = useState<null | string>(null);
    const [showQuickEdit, setShowQuickEdit] = useState(false);
    const [totalCount, setTotalCount] = useState<number>(0);
    const [page, setPage] = useState<number>(1);
    const limit = 10;
    const navigate = useNavigate();

    const [searchQuery, setSearchQuery] = useState('');
    const [appliedSearchQuery, setAppliedSearchQuery] = useState('');

    useEffect(() => {
        setLoading(true);
        const offset = (page - 1) * limit;
        
        const fetchProducts = async () => {
            try {
                let data;
                if (appliedSearchQuery) {
                    data = await productService.searchAdminProduct(appliedSearchQuery, limit);
                } else {
                    data = await productService.getproducts(offset, limit);
                }

                if (data === null) {
                    setProducts([]);
                    setTotalCount(0);
                } else if (data) {
                    if (Array.isArray(data)) {
                        setProducts(data);
                        setTotalCount(data.length);
                    } else if ('data' in data && data.data !== undefined) {
                        setProducts(data.data);
                        setTotalCount(data.count || 0);
                    } else {
                        setProducts([]);
                        setTotalCount(0);
                    }
                }
            } catch (error) {
                console.error("Hata oluştu:", error);
                setProducts([]);
                setTotalCount(0);
            } finally {
                setLoading(false);
            }
        };

        fetchProducts();
    }, [page, appliedSearchQuery]);


    const handleDelete = (id: string) => {
        if (window.confirm('Bu ürünü silmek istediğinize emin misiniz?')) {
            console.log("Silinecek ID:", id);
            productService.deleteProduct(id).then(() => {
                setProducts((prevProducts: GetCategoryWithProduct[]) => prevProducts.filter((p) => p.Product.product_id !== id));
                toast.success("Ürün başarıyla silindi.");
            }).catch((err) => {
                console.error("Silme işlemi başarısız:", err);
                toast.error("Ürün silinirken bir hata oluştu.");
            })
        }
    };

    const handleEdit = (id: string) => {
        navigate(`/admin/edit-product/${id}`);
    };

    const handleQuickEdit = (id: string) => {
        if (!showQuickEdit) {
            setShowQuickEdit(true);
            setSelectedProductId(id);
        }
    };

    if (loading) {
        return <Spinner message="Ürün listesi yükleniyor..." />
    }

    return (
        <Container maxWidth="xl" sx={{ mt: 2, mb: 4 }}>
            <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 3 }}>
                <Typography variant="h5" fontWeight="bold">
                    Ürün Listesi
                </Typography>
                <TextField
                    size="small"
                    placeholder="Ürün Ara..."
                    variant="outlined"
                    value={searchQuery}
                    onChange={(e) => setSearchQuery(e.target.value)}
                    onKeyDown={(e) => {
                        if (e.key === 'Enter') {
                            setAppliedSearchQuery(searchQuery.trim());
                            setPage(1);
                        }
                    }}
                    InputProps={{
                        endAdornment: (
                            <InputAdornment position="end">
                                <IconButton onClick={() => { setAppliedSearchQuery(searchQuery.trim()); setPage(1); }} edge="end">
                                    <SearchIcon />
                                </IconButton>
                            </InputAdornment>
                        )
                    }}
                    sx={{ width: 300, bgcolor: 'background.paper', borderRadius: 1 }}
                />
            </Box>

            <Stack spacing={2}>

                {products.length === 0 && (
                    <Box sx={{ textAlign: 'center', py: 4 }}>
                        <Typography color="text.secondary">Henüz ürün bulunmuyor.</Typography>
                    </Box>
                )}

                {products.map((productWithCategories: GetCategoryWithProduct) => (
                    <Paper
                        key={productWithCategories.Product.product_id}
                        sx={{
                            p: 2,
                            display: 'flex',
                            alignItems: 'center',
                            gap: 3,
                            transition: 'all 0.2s',
                            '&:hover': {
                                transform: 'translateY(-2px)',
                                boxShadow: 3
                            }
                        }}
                    >
                        {/* Ürün Görseli */}
                        <Avatar
                            src={productWithCategories.Product.thumbnail_image}
                            variant="rounded"
                            sx={{ width: 80, height: 80, bgcolor: 'grey.200' }}
                        >
                            {productWithCategories.Product.product_name.charAt(0)}
                        </Avatar>

                        {/* Ürün Bilgileri - Sol Kısım */}
                        <Box sx={{ flex: 1 }}>
                            <Typography variant="h6" fontWeight="bold">
                                {productWithCategories.Product.product_name}
                            </Typography>
                            <Box sx={{ display: 'flex', gap: 1, alignItems: 'center', mt: 0.5, flexWrap: 'wrap' }}>
                                <Typography variant="body2" color="text.secondary">
                                    Stok:
                                </Typography>
                                <Chip
                                    label={productWithCategories.Product.stock}
                                    size="small"
                                    color={productWithCategories.Product.stock < 10 ? 'error' : 'default'}
                                    variant="outlined"
                                />
                                {productWithCategories.Categories && productWithCategories.Categories.length > 0 && productWithCategories.Categories.map(cat => (
                                    <Chip
                                        key={cat.category_id}
                                        label={cat.category_name}
                                        size="small"
                                        color="primary"
                                        variant="filled"
                                        sx={{ height: 20, fontSize: '0.7rem' }}
                                    />
                                ))}
                            </Box>
                        </Box>

                        {/* Fiyat */}
                        <Typography variant="h6" color="primary.main" fontWeight="bold" sx={{ minWidth: '100px', textAlign: 'right' }}>
                            {productWithCategories.Product.price} ₺
                        </Typography>

                        {/* Aksiyon Butonları */}
                        <Stack direction="row" spacing={1}>
                            <Button
                                variant="outlined"
                                startIcon={<EditIcon />}
                                onClick={() => handleEdit(productWithCategories.Product.product_id)}
                                size="small"
                            >
                                Düzenle
                            </Button>

                            <Button
                                variant="outlined"
                                startIcon={<EditIcon />}
                                onClick={() => handleQuickEdit(productWithCategories.Product.product_id)}
                                size="small"
                            >
                                Hızlı Düzenle
                            </Button>
                            <Button
                                variant="outlined"
                                color="error"
                                startIcon={<DeleteIcon />}
                                onClick={() => handleDelete(productWithCategories.Product.product_id)}
                                size="small"
                            >
                                Sil
                            </Button>
                        </Stack>
                    </Paper>


                ))}



            </Stack>

            {totalCount > limit && (
                <Box sx={{ display: 'flex', justifyContent: 'center', mt: 4 }}>
                    <Pagination
                        count={Math.ceil(totalCount / limit)}
                        page={page}
                        onChange={(_event, value) => {
                            setPage(value);
                            window.scrollTo({ top: 0, behavior: 'smooth' });
                        }}
                        color="primary"
                    />
                </Box>
            )}

            {/* Dialogs */}

            {showQuickEdit && (
                <QuickEdit
                    open={showQuickEdit}
                    onClose={() => {
                        setShowQuickEdit(false);
                        setSelectedProductId(null);
                    }}
                    product={products.find(p => p.Product.product_id === selectedProductId)?.Product}
                />
            )}

        </Container>
    )
}

export default AdminProducts;