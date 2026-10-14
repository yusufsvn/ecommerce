import { Box, Button, Paper, Stack, TextField, Typography, Grid, MenuItem, Chip } from '@mui/material';
import { useState, useEffect } from 'react';
import type { Createproduct } from '../../models/ProductResponse';
import { productService } from '../../Services/productService';
import CloudUploadIcon from '@mui/icons-material/CloudUpload';
import DeleteIcon from '@mui/icons-material/Delete';
import { IconButton, LinearProgress } from '@mui/material';
import useFileUpload from '../../hooks/useFileUpload';
import { toast } from 'react-toastify';
import CreateProductFeature, { type ProductFeature } from './CreateProductFeature';

export default function CreateProduct() {
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
    const [formData, setFormData] = useState<Createproduct>({
        product_name: '',
        price: '0',
        description: '',
        product_category: [],
        stock: '0',
        images: [],
        item_type: 'PHYSICAL',
        technical_details: [],
        publish_state:'UNPUBLISHED'
    });



    // file upload process

    const { files, handleFileSelect, handleRemoveFile, handleThumbnailSelect, resetFiles, isChangedThumb } = useFileUpload();


    // tip dönüşümü metodu
    const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const { name, value } = e.target;
        setFormData((prev: Createproduct) => ({
            ...prev,
            [name]: name === 'price' ? parseFloat(value) : (name === 'stock' ? parseInt(value) : value)
        }));
    };

    const handleCategoryAdd = (e: React.ChangeEvent<HTMLInputElement>) => {
        const category = e.target.value;
        if (category && !formData.product_category.includes(category)) {
            setFormData((prev: Createproduct) => ({
                ...prev,
                product_category: [...prev.product_category, category]
            }));
        }
    };

    const handleCategoryDelete = (categoryToDelete: string) => {
        setFormData((prev: Createproduct) => ({
            ...prev,
            product_category: prev.product_category.filter(c => c !== categoryToDelete)
        }));
    };

    // submit 
    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        formData.images = files.map(file => file.file)
        console.log('Product Created:', formData);
        formData.technical_details = features

        try {
            await productService.createProduct(formData);
            toast.success("Ürün başarıyla oluşturuldu!");

            // Reset form
            setFormData({
                product_name: '',
                price: '0',
                description: '',
                product_category: [],
                stock: '0',
                images: [],
                item_type: 'PHYSICAL',
                technical_details: [],
                publish_state:'',
            });
            setFeatures([]);
            resetFiles();
        } catch (error) {
            console.error("Ürün oluşturma hatası:", error);
            // httpinterceptor will also show a generic error toast
        }
    };

    const [features, setFeatures] = useState<ProductFeature[]>([])

    return (
        <Box sx={{ maxWidth: 1200, mx: 'auto', mt: 4 }}>
            <Paper elevation={3} sx={{ p: 4 }}>
                <Typography variant="h5" component="h2" gutterBottom sx={{ mb: 4 }}>
                    Create New Product
                </Typography>
                <form onSubmit={handleSubmit}>
                    <Grid container spacing={4}>
                        {/* Left Column: Form Fields */}
                        <Grid size={{ xs: 12, md: 8 }}>
                            <Stack spacing={3}>
                                <TextField
                                    label="Ürün Adı"
                                    name="product_name"
                                    value={formData.product_name}
                                    onChange={handleChange}
                                    required
                                    fullWidth
                                />
                                <Stack direction={{ xs: 'column', sm: 'row' }} spacing={2}>
                                    <TextField
                                        label="Fiyat"
                                        name="price"
                                        type="number"
                                        value={formData.price}
                                        onChange={handleChange}
                                        required
                                        fullWidth
                                        slotProps={{ htmlInput: { min: 0, step: 0.01 } }}
                                    />
                                    <TextField
                                        label="Stok"
                                        name="stock"
                                        type="number"
                                        value={formData.stock}
                                        onChange={handleChange}
                                        fullWidth
                                        slotProps={{ htmlInput: { min: 0, step: 1 } }}
                                    />
                                </Stack>
                                <TextField
                                    select
                                    label="Kategori"
                                    name="product_category"
                                    value=""
                                    onChange={handleCategoryAdd}
                                    fullWidth
                                    required={formData.product_category.length === 0}
                                >
                                    {categories.map((option) => (
                                        <MenuItem key={option} value={option}>
                                            {option}
                                        </MenuItem>
                                    ))}
                                </TextField>
                                <Stack direction={"row"} spacing={1} flexWrap="wrap" useFlexGap>
                                    {formData.product_category.map(category => (
                                        <Chip
                                            key={category}
                                            label={category}
                                            onDelete={() => handleCategoryDelete(category)}
                                        />
                                    ))}
                                </Stack>

                                <Stack direction={"row"} spacing={2}> 
                                <TextField
                                    select
                                    label="Ürün Tipi"
                                    name="item_type"
                                    value={formData.item_type}
                                    onChange={handleChange}
                                    fullWidth
                                    required
                                >
                                    <MenuItem value="PHYSICAL">Fiziksel</MenuItem>
                                    <MenuItem value="VIRTUAL">Sanal</MenuItem>
                                </TextField>
                                <TextField
                                    select
                                    label="Yayınla"
                                    name="publish_state"
                                    value={formData.publish_state}
                                    onChange={handleChange}
                                    fullWidth
                                    required
                                >
                                    <MenuItem value="PUBLISHED">Yayınla</MenuItem>
                                    <MenuItem value="UNPUBLISHED">Yayından Kaldır</MenuItem>
                                </TextField>
                                </Stack>

                                <CreateProductFeature features={features} onChange={setFeatures}></CreateProductFeature>

                                <TextField
                                    label="Description"
                                    name="description"
                                    value={formData.description}
                                    onChange={handleChange}
                                    multiline
                                    rows={4}
                                    fullWidth
                                />
                            </Stack>
                        </Grid>

                        {/* Right Column: Actions */}
                        <Grid size={{ xs: 12, md: 4 }}>
                            <Stack spacing={2} sx={{ height: '100%', justifyContent: 'flex-start' }}>
                                {/* Thumbnail Preview Box */}
                                <Box sx={{
                                    width: '100%',
                                    height: 200,
                                    border: '2px dashed',
                                    borderColor: 'divider',
                                    borderRadius: 2,
                                    display: 'flex',
                                    justifyContent: 'center',
                                    alignItems: 'center',
                                    overflow: 'hidden',
                                    mb: 2,
                                    bgcolor: 'background.paper'
                                }}>
                                    {isChangedThumb && files.length > 0 ? (
                                        <Box
                                            component="img"
                                            src={files[0].preview}
                                            alt="Kapak Fotoğrafı"
                                            sx={{ width: '100%', height: '100%', objectFit: 'cover' }}
                                        />
                                    ) : (
                                        <Typography variant="body2" color="text.secondary">
                                            Kapak Fotoğrafı Yüklenmedi
                                        </Typography>
                                    )}
                                </Box>

                                <Button
                                    component="label"
                                    variant="outlined"
                                    startIcon={<CloudUploadIcon />}
                                    fullWidth
                                    sx={{ height: 56 }}
                                >
                                    Kapak Fotoğrafı Yükle
                                    <input
                                        type="file"
                                        hidden
                                        accept="image/*"
                                        onChange={handleThumbnailSelect}
                                    />
                                </Button>

                                <Button
                                    component="label"
                                    variant="outlined"
                                    startIcon={<CloudUploadIcon />}
                                    fullWidth
                                    sx={{ height: 56 }}
                                >
                                    Ürün Resimleri Yükle
                                    <input
                                        type="file"
                                        hidden
                                        multiple
                                        accept="image/*"
                                        onChange={handleFileSelect}
                                    />
                                </Button>

                                <Button
                                    variant="contained"
                                    color="primary"
                                    type="submit"
                                    size="large"
                                    fullWidth
                                    sx={{ height: 56 }}
                                >
                                    Create Product
                                </Button>

                                {/* Image Previews */}
                                {files.length > (isChangedThumb ? 1 : 0) && (
                                    <Box sx={{ mb: 2 }}>
                                        <Grid container spacing={2}>
                                            {(isChangedThumb ? files.slice(1) : files).map((file) => (
                                                <Grid size={{ xs: 12 }} key={file.id}>
                                                    <Paper elevation={2} sx={{ p: 2, display: 'flex', alignItems: 'center', gap: 2 }}>
                                                        <Box
                                                            component="img"
                                                            src={file.preview}
                                                            alt="Preview"
                                                            sx={{ width: 60, height: 60, objectFit: 'cover', borderRadius: 1 }}
                                                        />
                                                        <Box sx={{ flexGrow: 1 }}>
                                                            <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 0.5 }}>
                                                                <Typography variant="body2" noWrap sx={{ maxWidth: 150 }}>
                                                                    {file.file.name}
                                                                </Typography>
                                                                {file.status === 'completed' && (
                                                                    <Typography variant="caption" color="success.main">
                                                                        Tamamlandı
                                                                    </Typography>
                                                                )}
                                                            </Box>
                                                            <LinearProgress
                                                                variant="determinate"
                                                                value={file.progress}
                                                                color={file.status === 'completed' ? 'success' : 'primary'}
                                                            />
                                                        </Box>
                                                        <IconButton
                                                            size="small"
                                                            color="error"
                                                            onClick={() => handleRemoveFile(file.id)}
                                                        >
                                                            <DeleteIcon />
                                                        </IconButton>
                                                    </Paper>
                                                </Grid>
                                            ))}
                                        </Grid>
                                    </Box>
                                )}


                            </Stack>
                        </Grid>
                    </Grid>
                </form>
            </Paper>
        </Box>
    );
}
