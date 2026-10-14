import { Box, Stack, TextField, Button, ImageList, ImageListItem, ImageListItemBar, IconButton, Typography, Paper, Container, MenuItem } from '@mui/material';
import { useState, useEffect } from 'react';
import type { FullProductResponse, GetCategory, ProductImage } from '../../../models/ProductResponse';
import { productService } from '../../../Services/productService';
import DeleteIcon from '@mui/icons-material/Delete';
import CloudUploadIcon from '@mui/icons-material/CloudUpload';
import useFileUpload from '../../../hooks/useFileUpload';
import { toast } from 'react-toastify';
import Chip from '@mui/material/Chip';
import EditIcon from '@mui/icons-material/Edit';
import CreateProductFeature, { type ProductFeature } from '../CreateProductFeature';
import {
    publishStatusOptions,
    itemTypeOptions,
    mapTurkishToEnglishPublishState,
    mapEnglishToTurkishPublishState,
    mapTurkishToEnglishItemType,
    mapEnglishToTurkishItemType
} from '../../../utils/statusMapper';

import { useParams, useNavigate } from 'react-router';

function EditProduct() {
    const { id: productId } = useParams<{ id: string }>();
    const navigate = useNavigate();
    const [product, setProduct] = useState<FullProductResponse | null>(null);
    const [categories, setCategories] = useState<string[]>([]);

    useEffect(() => {
        productService.getCategories().then((data) => {
            if (data && data.length > 0) {
                setCategories(data.map(c => c.category_name));
            }
        }).catch(err => {
            console.error("Kategoriler yüklenirken hata:", err);
        });
    }, []);



    const [formData, setFormData] = useState({
        product_id: "",
        product_name: '',
        stock: 0,
        price: 0,
        description: '',
        product_code: '',
        images: [] as ProductImage[],
        categories: [] as GetCategory[],
        publishState: 'Yayından Kaldır',
        itemType: 'Fiziksel'
    });

    const [deletedImages, setDeletedImages] = useState<string[]>([]);

    // seçilen fotoğraflar files
    const { files, isChangedThumb, handleFileSelect, handleRemoveFile, handleThumbnailSelect } = useFileUpload();


    // Kategoriler için state
    const [deletedCategories, setDeletedCategories] = useState<string[]>([]);
    const [newCategories, setNewCategories] = useState<string[]>([]);
    const [features, setFeatures] = useState<ProductFeature[]>([]);


    useEffect(() => {

        // Fetch the product by ID
        if (productId) {
            productService.getProductByIdForAdmin(productId).then((data) => {
                setProduct(data);

                if (data) {
                    const rawPublishState = data.product_publish_state || data.ProductPublishState || '';
                    const rawItemType = data.item_type || data.ItemType || '';

                    setFormData({
                        product_id: data.product_id || data.ProductID || '',
                        product_name: data.product_name || data.ProductName || '',
                        stock: data.stock ?? data.Stock ?? 0,
                        price: data.price ?? data.Price ?? 0,
                        description: data.description || data.Description || '',
                        product_code: data.product_code || data.ProductCode || '',
                        images: data.images || data.Images || [],
                        categories: data.categories || data.Categories || [],
                        publishState: mapEnglishToTurkishPublishState(rawPublishState),
                        itemType: mapEnglishToTurkishItemType(rawItemType)
                    });
                    const techDetails = data.technical_details || data.TechnicalDetails;
                    if (techDetails) {
                        const loadedFeatures = Object.entries(techDetails).map(([key, value]) => ({
                            key,
                            value
                        }));
                        setFeatures(loadedFeatures);
                    }
                }
            })
        }
    }, [productId]);

    if (!product) {
        return null;
    }

    const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const { name, value } = e.target;
        setFormData(prev => ({
            ...prev,
            [name]: name === 'stock' || name === 'price' ? Number(value) : value
        }));
    };

    const handleDeleteImage = (imageId: string) => {
        const isConfirmed = window.confirm("Bu resmi silmek istediğinize emin misiniz?");
        if (isConfirmed) {
            if (imageId.length < 12) {
                handleRemoveFile(imageId);
            }
            setDeletedImages(prev => [...prev, imageId]);

            // bildirim
            toast.info("Resim silinecekler listesine eklendi.");
        } else {
            // 3. İptal ederse hiçbir şey yapma
            console.log("İşlem iptal edildi.");
        }
    };

    const handleSubmit = () => {
        const currentProductId = product.product_id || product.ProductID || productId || '';

        // call the add image request
        if (files.length > 0) {
            productService.addProductImage(files.map(file => file.file), currentProductId, isChangedThumb).then((data) => {
                if (data) {
                    toast.success("Yeni resimler eklendi.");
                }
            })
        }

        const detailsMap: Record<string, string> = {};
        features.forEach(f => {
            detailsMap[f.key] = f.value;
        });

        const publishStateEnglish = mapTurkishToEnglishPublishState(formData.publishState);
        const itemTypeEnglish = mapTurkishToEnglishItemType(formData.itemType);

        productService.updateProduct({
            ProductID: currentProductId,
            ProductName: formData.product_name,
            Stock: formData.stock,
            Price: formData.price,
            Categories: {
                NewCategories: newCategories,
                DeletedCategories: deletedCategories
            },
            Description: formData.description,
            DeletedImages: deletedImages,
            TechnicalDetails: detailsMap,
            ProductPublishState: publishStateEnglish,
            ItemType: itemTypeEnglish
        }).then((data) => {
            if (data) {
                toast.success("Ürün başarıyla güncellendi.");
                navigate('/admin/products');
            }
        })
    }

    const handleCancel = () => {
        navigate('/admin/products');
    }

    const handleCategoryAdd = (e: React.ChangeEvent<HTMLInputElement>) => {
        const category = e.target.value;
        if (category && !formData.categories.some(c => c.category_name === category)) {
            setFormData(prev => ({
                ...prev,
                categories: [...prev.categories, { category_id: '', category_name: category }]
            }));

            // Eğer daha önce silinmişler arasındaysa silinenlerden çıkar
            if (deletedCategories.includes(category)) {
                setDeletedCategories(prev => prev.filter(c => c !== category));
            } else {
                setNewCategories(prev => [...prev, category]);
            }
        }
    };

    const handleCategoryDelete = (categoryToDelete: string) => {
        setFormData(prev => ({
            ...prev,
            categories: prev.categories.filter(c => c.category_name !== categoryToDelete)
        }));

        // Eğer yeni eklenmiş bir kategoriyse, newCategories listesinden çıkar, değilse deletedCategories'e ekle
        if (newCategories.includes(categoryToDelete)) {
            setNewCategories(prev => prev.filter(c => c !== categoryToDelete));
        } else {
            setDeletedCategories(prev => [...prev, categoryToDelete]);
        }
    };


    return (
        <Container maxWidth="lg" sx={{ mt: 4, mb: 4 }}>
            <Paper sx={{ p: 4 }}>
                <Typography variant="h5" gutterBottom sx={{ mb: 3, fontWeight: 'bold' }}>Detaylı Düzenle</Typography>
                <Box sx={{ display: 'flex', flexDirection: { xs: 'column', md: 'row' }, gap: 4, mt: 1 }}>
                    {/* Left Column - Images */}
                    <Box sx={{ width: { xs: '100%', md: '45%' }, display: 'flex', flexDirection: 'column', gap: 2 }}>
                        {/* Kapak Fotoğrafı */}
                        {formData.images.length > 0 || (isChangedThumb && files.length > 0) ? (
                            <Box sx={{ position: 'relative', width: '100%', maxWidth: '250px', mx: 'auto' }}>
                                <img
                                    src={isChangedThumb && files.length > 0 ? files[0].preview : `${formData.images[0]?.path}?w=400&fit=crop&auto=format`}
                                    alt="Kapak Fotoğrafı"
                                    style={{ width: '100%', height: 'auto', borderRadius: '8px', objectFit: 'cover', display: 'block' }}
                                />
                                <IconButton
                                    component="label"
                                    size="small"
                                    sx={{ position: 'absolute', top: 8, right: 8, bgcolor: 'rgba(255,255,255,0.8)', '&:hover': { bgcolor: 'rgba(255,255,255,1)' } }}
                                >
                                    <input
                                        type="file"
                                        hidden
                                        multiple
                                        onChange={(e) => handleThumbnailSelect(e)}
                                    />
                                    <EditIcon color="error" fontSize="small" />
                                </IconButton>
                                <Typography variant="caption" sx={{ display: 'block', textAlign: 'center', mt: 1, color: 'text.secondary', fontWeight: 'bold' }}>
                                    Kapak Fotoğrafı
                                </Typography>
                            </Box>
                        ) : (
                            <Box sx={{ width: '100%', maxWidth: '250px', mx: 'auto', aspectRatio: '1', bgcolor: 'grey.200', display: 'flex', alignItems: 'center', justifyContent: 'center', borderRadius: '8px' }}>
                                <Typography color="text.secondary">Kapak Fotoğrafı Yok</Typography>
                            </Box>
                        )}

                        {/* Diğer Fotoğraflar */}
                        <Box>
                            <Typography variant="subtitle2" sx={{ mb: 1, mt: 1 }}>Diğer Fotoğraflar</Typography>
                            <ImageList cols={3} gap={8} rowHeight={80} sx={{ m: 0 }}>
                                {formData.images.slice(1).map((item) => (
                                    <ImageListItem key={item.image_id}>
                                        <img
                                            src={`${item.path}?w=80&h=80&fit=crop&auto=format`}
                                            alt={item.image_id}
                                            loading="lazy"
                                            style={{ borderRadius: '4px', height: '100%', objectFit: 'cover' }}
                                        />
                                        <ImageListItemBar
                                            actionIcon={
                                                <IconButton
                                                    size="small"
                                                    sx={{ color: 'rgba(255, 0, 0, 0.87)' }}
                                                    onClick={() => handleDeleteImage(item.image_id)}
                                                >
                                                    <DeleteIcon fontSize="small" />
                                                </IconButton>
                                            }
                                        />
                                    </ImageListItem>
                                ))}
                                {(isChangedThumb && files.length > 0 ? files.filter(f => f.id !== files[0].id) : files).map((file) => (
                                    <ImageListItem key={file.id}>
                                        <img
                                            src={file.preview}
                                            alt="New upload"
                                            loading="lazy"
                                            style={{ opacity: file.status === 'uploading' ? 0.7 : 1, width: '100%', height: '100%', objectFit: 'cover', borderRadius: '4px' }}
                                        />
                                        <ImageListItemBar
                                            actionIcon={
                                                <IconButton
                                                    size="small"
                                                    sx={{ color: 'rgba(255, 0, 0, 0.87)' }}
                                                    onClick={() => handleRemoveFile(file.id)}
                                                >
                                                    <DeleteIcon fontSize="small" />
                                                </IconButton>
                                            }
                                        />
                                    </ImageListItem>
                                ))}
                            </ImageList>
                        </Box>
                    </Box>

                    {/* Right Column - Form Fields */}
                    <Box sx={{ width: { xs: '100%', md: '55%' }, display: 'flex', flexDirection: 'column', gap: 2 }}>
                        {formData.product_code && (
                            <TextField
                                label="Ürün Kodu"
                                name="product_code"
                                fullWidth
                                value={formData.product_code}
                                disabled
                            />
                        )}
                        <TextField
                            label="Ürün Adı"
                            name="product_name"
                            fullWidth
                            value={formData.product_name}
                            onChange={handleChange}
                        />
                        <Stack direction="row" spacing={2}>
                            <TextField
                                label="Stok"
                                name="stock"
                                type="number"
                                fullWidth
                                value={formData.stock}
                                onChange={handleChange}
                            />
                            <TextField
                                label="Fiyat"
                                name="price"
                                type="number"
                                fullWidth
                                value={formData.price}
                                onChange={handleChange}
                            />
                        </Stack>
                        <Stack direction="row" spacing={2}>
                            <TextField
                                select
                                label="Yayın Durumu"
                                name="publishState"
                                fullWidth
                                value={formData.publishState}
                                onChange={handleChange}
                            >
                                {publishStatusOptions.map((option) => (
                                    <MenuItem key={option} value={option}>
                                        {option}
                                    </MenuItem>
                                ))}
                            </TextField>
                            <TextField
                                select
                                label="Ürün Tipi"
                                name="itemType"
                                fullWidth
                                value={formData.itemType}
                                onChange={handleChange}
                            >
                                {itemTypeOptions.map((option) => (
                                    <MenuItem key={option} value={option}>
                                        {option}
                                    </MenuItem>
                                ))}
                            </TextField>
                        </Stack>
                        <TextField
                            label="Açıklama"
                            name="description"
                            fullWidth
                            multiline
                            rows={6}
                            value={formData.description}
                            onChange={handleChange}
                        />


                        <Typography variant="h6">Categories</Typography>

                        <TextField
                            select
                            label="Kategori Ekle"
                            value=""
                            onChange={handleCategoryAdd}
                            fullWidth
                        >
                            {categories.map((option) => (
                                <MenuItem key={option} value={option}>
                                    {option}
                                </MenuItem>
                            ))}
                        </TextField>

                        <Stack direction={"row"} spacing={1} flexWrap="wrap" useFlexGap>
                            {formData.categories.map(category => (
                                <Chip
                                    key={category.category_name}
                                    label={category.category_name}
                                    onDelete={() => handleCategoryDelete(category.category_name)}
                                />
                            ))}
                        </Stack>

                        <CreateProductFeature features={features} onChange={setFeatures}></CreateProductFeature>

                        <Button
                            component="label"
                            variant="outlined"
                            startIcon={<CloudUploadIcon />}
                            sx={{ alignSelf: 'flex-start', mt: 1 }}
                        >
                            Daha Fazla Resim Yükle
                            <input
                                type="file"
                                hidden
                                multiple
                                onChange={(e) => handleFileSelect(e)}
                            />
                        </Button>
                    </Box>
                </Box>
                <Box sx={{ mt: 4, display: 'flex', justifyContent: 'flex-end', gap: 2 }}>
                    <Button onClick={handleCancel}>İptal</Button>
                    <Button variant="contained" onClick={handleSubmit}>Kaydet</Button>
                </Box>
            </Paper>
        </Container>
    )
}

export default EditProduct;
