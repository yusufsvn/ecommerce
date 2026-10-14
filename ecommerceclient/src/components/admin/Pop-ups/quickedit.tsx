import { Dialog, DialogTitle, DialogContent, Stack, TextField, DialogActions, Button } from '@mui/material';
import { useState, useEffect } from 'react';
import type { GetProducts } from '../../../models/ProductResponse';
import { productService } from '../../../Services/productService';
import { toast } from 'react-toastify';

interface QuickEditProps {
    open: boolean;
    onClose: () => void;
    product: GetProducts | undefined;
}

function QuickEdit({ open, onClose, product }: QuickEditProps) {
    const [formData, setFormData] = useState({
        ProductName: '',
        Stock: 0,
        Price: 0
    });

    useEffect(() => {
        if (product) {
            setFormData({
                ProductName: product.product_name,
                Stock: product.stock,
                Price: product.price
            });
        }
    }, [product]);

    if (!product) {
        return null;
    }

    const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const { name, value } = e.target;
        setFormData(prev => ({
            ...prev,
            // stock ve price ın number turunde olması gerektiğinden
            [name]: name === 'Stock' || name === 'Price' ? Number(value) : value
        }));
    };

    const handleSubmit = () => {
        productService.patchProduct(
            {
                ProductID: product.product_id,
                ProductName: formData.ProductName,
                Stock: formData.Stock,
                Price: formData.Price
            }
        ).then((data) => {
            if (data) {
                toast.success("Ürün başarıyla güncellendi.");
                onClose();
            }
        })
    }

    return (
        <Dialog open={open} onClose={onClose} maxWidth="sm" fullWidth>
            <DialogTitle>Hızlı Düzenle</DialogTitle>
            <DialogContent>
                <Stack spacing={2} sx={{ mt: 1 }}>
                    <TextField
                        label="Ürün Adı"
                        name="ProductName"
                        fullWidth
                        value={formData.ProductName}
                        onChange={handleChange}
                    />
                    <TextField
                        label="Stok"
                        name="Stock"
                        type="number"
                        value={formData.Stock}
                        onChange={handleChange}
                        fullWidth
                    />
                    <TextField
                        label="Fiyat"
                        name="Price"
                        type="number"
                        value={formData.Price}
                        onChange={handleChange}
                        fullWidth
                    />
                </Stack>
            </DialogContent>
            <DialogActions>
                <Button onClick={onClose}>İptal</Button>
                <Button variant="contained" onClick={handleSubmit}>Kaydet</Button>
            </DialogActions>
        </Dialog>
    )
}

export default QuickEdit;