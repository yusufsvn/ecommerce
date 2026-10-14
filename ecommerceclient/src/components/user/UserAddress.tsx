import React, { useEffect, useState } from 'react';
import { Container, Paper, Typography, Box, TextField, Button, Grid, Stack, Card, CardContent, CardActions } from '@mui/material';
import { useNavigate } from 'react-router';
import { toast } from 'react-toastify';
import { userService } from '../../Services/userService';
import type { addressRequest } from '../../models/addressmodel';
import ArrowBackIcon from '@mui/icons-material/ArrowBack';
import MapIcon from '@mui/icons-material/Map';
import EditIcon from '@mui/icons-material/Edit';
import DeleteIcon from '@mui/icons-material/Delete';
import AddIcon from '@mui/icons-material/Add';
import HomeIcon from '@mui/icons-material/Home';

const UserAddress = () => {
    const navigate = useNavigate();

    const [addresses, setAddresses] = useState<addressRequest[]>([]);
    const [isFormOpen, setIsFormOpen] = useState(false);
    const [editingIndex, setEditingIndex] = useState<number | null>(null);

    const [addressData, setAddressData] = useState<addressRequest>({
        AddressName: '',
        AddressStr: '',
        City: '',
        District: '',
        Country: '',
        ZipCode: '',
    });

    useEffect(() => {
        // get address data
        userService.getAllUserAddress().then((res) => {
            if (res && res.Addresses) {
                setAddresses(res.Addresses);
            }
        }).catch(err => {
            console.error("Adres alınırken bir hata oluştu", err);
        });
    }, []);

    const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const { name, value } = e.target;
        setAddressData(prev => ({
            ...prev,
            [name]: value
        }));
    };

    const handleCreate = async (e: React.FormEvent) => {
        e.preventDefault();

        userService.addAddress(addressData).then((res) => {
            if (res === 200) {
                toast.success('Yeni adres başarıyla eklendi!');

                userService.getAllUserAddress().then((res) => {
                    if (res && res.Addresses) {
                        setAddresses(res.Addresses);
                    } else {
                        setAddresses(prev => [...prev, addressData]);
                    }
                }).catch(err => {
                    console.error("Adres alınırken bir hata oluştu", err);
                    setAddresses(prev => [...prev, addressData]);
                });

                // Formu kapat ve state'leri sıfırla
                setIsFormOpen(false);
                setAddressData({ AddressName: '', AddressStr: '', City: '', District: '', Country: '', ZipCode: '' });
            } else {
                toast.error('Yeni adres eklenirken bir hata oluştu!');
            }
        }).catch(() => {
            toast.error('Sunucu ile bağlantı kurulamadı!');
        });
    };

    const handleUpdate = async (e: React.FormEvent) => {
        e.preventDefault();
        if (editingIndex === null) return;

        userService.updateAddress(addressData).then((res) => {
            if (res === 200) {
                toast.success('Adres bilgileriniz başarıyla güncellendi!');

                const newAddresses = [...addresses];
                newAddresses[editingIndex] = addressData;
                setAddresses(newAddresses);

                // Formu kapat ve state'leri sıfırla
                setIsFormOpen(false);
                setEditingIndex(null);
                setAddressData({ AddressName: '', AddressStr: '', City: '', District: '', Country: '', ZipCode: '' });
            } else {
                toast.error('Adres bilgileriniz güncellenirken bir hata oluştu!');
            }
        }).catch(() => {
            toast.error('Sunucu ile bağlantı kurulamadı!');
        });
    };


    const handleEdit = (index: number) => {
        setAddressData(addresses[index]);
        setEditingIndex(index);
        setIsFormOpen(true);
    };

    const handleDelete = async (index: number) => {
        await userService.deleteAddress(addresses[index].AddressId ?? "");
        const res = await userService.getAllUserAddress();
        if (res && res.Addresses) {
            setAddresses(res.Addresses);
            toast.success('Adres başarıyla silindi!');
        } else {
            toast.error('Adres silinirken bir hata oluştu!');
        }
    };

    const handleAddNew = async () => {
        setAddressData({ AddressName: '', AddressStr: '', City: '', District: '', Country: '', ZipCode: '' });
        setEditingIndex(null);
        setIsFormOpen(true);
    };

    const handleCancel = () => {
        setIsFormOpen(false);
        setEditingIndex(null);
    };

    return (
        <Container component="main" maxWidth="md" sx={{ py: { xs: 4, md: 8 } }}>
            <Box sx={{ mb: 4, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                <Button
                    startIcon={<ArrowBackIcon />}
                    onClick={() => navigate('/profile')}
                    sx={{ color: 'text.secondary', '&:hover': { bgcolor: 'grey.100' }, borderRadius: 8, px: 2 }}
                >
                    Profile Dön
                </Button>
                {!isFormOpen && (
                    <Button
                        startIcon={<AddIcon />}
                        variant="contained"
                        color="primary"
                        onClick={handleAddNew}
                        sx={{ borderRadius: 8, px: 3, py: 1, boxShadow: 2 }}
                    >
                        Yeni Adres Ekle
                    </Button>
                )}
            </Box>

            <Paper elevation={0} sx={{ p: { xs: 3, md: 6 }, borderRadius: 4, border: '1px solid', borderColor: 'divider' }}>
                <Box sx={{ display: 'flex', alignItems: 'center', mb: 4, justifyContent: 'center', flexDirection: 'column' }}>
                    <Box sx={{ p: 2, borderRadius: '50%', bgcolor: 'primary.50', mb: 2 }}>
                        <MapIcon color="primary" sx={{ fontSize: 40 }} />
                    </Box>
                    <Typography component="h1" variant="h4" fontWeight="bold" align="center">
                        Adres Bilgilerim
                    </Typography>
                    <Typography variant="body1" color="text.secondary" align="center" sx={{ mt: 1 }}>
                        Siparişlerinizin sorunsuz teslimatı için adres bilgilerinizi güncel tutun.
                    </Typography>
                </Box>

                {!isFormOpen ? (
                    <Stack spacing={3} sx={{ mt: 4 }}>
                        {addresses.length === 0 ? (
                            <Box sx={{ p: 5, textAlign: 'center', bgcolor: 'grey.50', borderRadius: 3, border: '1px dashed', borderColor: 'grey.300' }}>
                                <Typography variant="h6" color="text.secondary" gutterBottom>
                                    Kayıtlı adresiniz bulunmamaktadır
                                </Typography>
                                <Typography variant="body2" color="text.secondary" sx={{ mb: 3 }}>
                                    Hemen yeni bir adres ekleyerek alışverişe devam edebilirsiniz.
                                </Typography>
                                <Button variant="outlined" startIcon={<AddIcon />} onClick={handleAddNew} sx={{ borderRadius: 8 }}>
                                    Adres Ekle
                                </Button>
                            </Box>
                        ) : (
                            addresses.map((address, index) => (
                                <Card key={index} variant="outlined" sx={{ borderRadius: 3, borderColor: 'grey.200', transition: 'box-shadow 0.3s', '&:hover': { boxShadow: 3 } }}>
                                    <CardContent sx={{ display: 'flex', alignItems: 'flex-start', gap: 2, pb: 1 }}>
                                        <Box sx={{ mt: 0.5, color: 'primary.main' }}>
                                            <HomeIcon />
                                        </Box>
                                        <Box sx={{ flex: 1 }}>
                                            <Typography variant="h6" gutterBottom fontWeight="bold" color="text.primary">
                                                {address.AddressName || 'Kayıtlı Adres'}
                                            </Typography>
                                            <Typography variant="subtitle2" color="text.secondary" sx={{ mb: 1 }}>
                                                {address.City} / {address.District}
                                            </Typography>
                                            <Typography variant="body1" color="text.secondary" paragraph sx={{ mb: 1 }}>
                                                {address.AddressStr}
                                            </Typography>
                                            <Typography variant="body2" color="text.secondary" fontWeight="medium">
                                                {address.ZipCode} - {address.Country}
                                            </Typography>
                                        </Box>
                                    </CardContent>
                                    <CardActions sx={{ justifyContent: 'flex-end', px: 3, pb: 2, pt: 0 }}>
                                        <Button size="small" variant="text" startIcon={<EditIcon />} onClick={() => handleEdit(index)} sx={{ borderRadius: 2, fontWeight: 'bold' }}>
                                            Düzenle
                                        </Button>
                                        <Button size="small" variant="text" color="error" startIcon={<DeleteIcon />} onClick={() => handleDelete(index)} sx={{ borderRadius: 2, fontWeight: 'bold' }}>
                                            Sil
                                        </Button>
                                    </CardActions>
                                </Card>
                            ))
                        )}
                    </Stack>
                ) : (
                    <Box component="form" onSubmit={editingIndex !== null ? handleUpdate : handleCreate} sx={{ mt: 4, p: 3, bgcolor: 'grey.50', borderRadius: 3, border: '1px solid', borderColor: 'grey.200' }}>
                        <Typography variant="h6" fontWeight="bold" sx={{ mb: 3 }}>
                            {editingIndex !== null ? 'Adresi Düzenle' : 'Yeni Adres Ekle'}
                        </Typography>
                        <Grid container spacing={3}>
                            <Grid size={{ xs: 12 }}>
                                <TextField
                                    required
                                    fullWidth
                                    id="addressName"
                                    label="Adres Başlığı (Örn: Ev, İş)"
                                    name="AddressName"
                                    value={addressData.AddressName || ''}
                                    onChange={handleChange}
                                    sx={{ '& .MuiOutlinedInput-root': { borderRadius: 2, bgcolor: 'background.paper' } }}
                                />
                            </Grid>
                            <Grid size={{ xs: 12 }}>
                                <TextField
                                    required
                                    fullWidth
                                    id="address"
                                    label="Açık Adres"
                                    name="AddressStr"
                                    value={addressData.AddressStr}
                                    onChange={handleChange}
                                    multiline
                                    rows={3}
                                    sx={{ '& .MuiOutlinedInput-root': { borderRadius: 2, bgcolor: 'background.paper' } }}
                                />
                            </Grid>
                            <Grid size={{ xs: 12, sm: 6 }}>
                                <TextField
                                    required
                                    fullWidth
                                    id="city"
                                    label="İl"
                                    name="City"
                                    value={addressData.City}
                                    onChange={handleChange}
                                    sx={{ '& .MuiOutlinedInput-root': { borderRadius: 2, bgcolor: 'background.paper' } }}
                                />
                            </Grid>
                            <Grid size={{ xs: 12, sm: 6 }}>
                                <TextField
                                    required
                                    fullWidth
                                    id="district"
                                    label="İlçe"
                                    name="District"
                                    value={addressData.District}
                                    onChange={handleChange}
                                    sx={{ '& .MuiOutlinedInput-root': { borderRadius: 2, bgcolor: 'background.paper' } }}
                                />
                            </Grid>
                            <Grid size={{ xs: 12, sm: 6 }}>
                                <TextField
                                    required
                                    fullWidth
                                    id="country"
                                    label="Ülke"
                                    name="Country"
                                    value={addressData.Country}
                                    onChange={handleChange}
                                    sx={{ '& .MuiOutlinedInput-root': { borderRadius: 2, bgcolor: 'background.paper' } }}
                                />
                            </Grid>
                            <Grid size={{ xs: 12, sm: 6 }}>
                                <TextField
                                    required
                                    fullWidth
                                    id="zip_code"
                                    label="Posta Kodu"
                                    name="ZipCode"
                                    value={addressData.ZipCode}
                                    onChange={handleChange}
                                    sx={{ '& .MuiOutlinedInput-root': { borderRadius: 2, bgcolor: 'background.paper' } }}
                                />
                            </Grid>
                        </Grid>

                        <Box sx={{ mt: 4, display: 'flex', justifyContent: 'flex-end', gap: 2 }}>
                            <Button
                                onClick={handleCancel}
                                size="large"
                                variant="outlined"
                                color="inherit"
                                sx={{ px: 4, py: 1.2, borderRadius: 8, fontWeight: 'bold' }}
                            >
                                İptal
                            </Button>
                            <Button
                                type="submit"
                                size="large"
                                variant="contained"
                                color="primary"
                                sx={{ px: 5, py: 1.2, borderRadius: 8, fontWeight: 'bold' }}
                            >
                                Kaydet
                            </Button>
                        </Box>
                    </Box>
                )}
            </Paper>
        </Container>
    );
};

export default UserAddress;
