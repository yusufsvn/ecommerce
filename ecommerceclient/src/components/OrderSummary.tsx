import { Box, Button, Container, Divider, TextField, Typography, Snackbar, Alert, FormControlLabel, Checkbox, Paper, Grid, Stepper, Step, StepLabel, Radio, RadioGroup, FormControl } from "@mui/material";
import { useState, useEffect } from "react";
import { useNavigate } from "react-router";
import type { cartItem } from "../models/cartmodel";
import type { CreateOrderRequest } from "../models/ordermodel";
import { orderService } from "../Services/orderService";
import { userService } from "../Services/userService";
import type { addressRequest } from "../models/addressmodel";
import SecurityIcon from '@mui/icons-material/Security';

interface OrderSummaryProps {
    cartItems: cartItem[];
    totalPrice: number;
    clearCart: () => void;
}

function OrderSummary({ totalPrice, clearCart }: OrderSummaryProps) {
    const navigate = useNavigate();
    const [loading, setLoading] = useState(false);
    const [activeStep, setActiveStep] = useState(0);
    const [notification, setNotification] = useState<{ open: boolean; message: string; severity: 'success' | 'error' }>({
        open: false,
        message: '',
        severity: 'success'
    });

    const [savedAddresses, setSavedAddresses] = useState<addressRequest[]>([]);
    const [selectedAddressIndex, setSelectedAddressIndex] = useState<number>(-1);
    const [saveNewAddress, setSaveNewAddress] = useState<boolean>(false);

    const [useSameAddressForBilling, setUseSameAddressForBilling] = useState<boolean>(true);
    const [selectedBillingAddressIndex, setSelectedBillingAddressIndex] = useState<number>(-1);

    const steps = ['Kimlik', 'Teslimat', 'Fatura', 'Ödeme'];

    // Form state
    const [formData, setFormData] = useState({
        firstName: "",
        lastName: "",
        email: "",
        phone: "",
        address: "",
        city: "",
        district: "",
        zipCode: "",
        cardName: "",
        cardNumber: "",
        expireMonth: "",
        expireYear: "",
        cvv: "",
        tcNo: "",
        billingAddress: "",
        billingCity: "",
        billingDistrict: "",
        billingZipCode: "",
        use3DSecure: false
    });

    useEffect(() => {
        userService.getAllUserAddress().then((res) => {
            if (res && res.Addresses && res.Addresses.length > 0) {

                setSavedAddresses(res.Addresses);
                setSelectedAddressIndex(0);
                setSelectedBillingAddressIndex(0);
                setFormData(prevData => ({
                    ...prevData,
                    address: res.Addresses[0].AddressStr,
                    city: res.Addresses[0].City,
                    district: res.Addresses[0].District,
                    zipCode: res.Addresses[0].ZipCode
                }));
            } else if (savedAddresses.length > 0) {
                setSelectedAddressIndex(0);
                setSelectedBillingAddressIndex(0);
                setFormData(prevData => ({
                    ...prevData,
                    address: savedAddresses[0].AddressStr,
                    city: savedAddresses[0].City,
                    district: savedAddresses[0].District,
                    zipCode: savedAddresses[0].ZipCode
                }));
            } else {
                setSelectedAddressIndex(-1);
                setSelectedBillingAddressIndex(-1);
            }
        }).catch(err => {
            console.log(err);
            if (savedAddresses.length > 0) {
                setSelectedAddressIndex(0);
                setSelectedBillingAddressIndex(0);
                setFormData(prevData => ({
                    ...prevData,
                    address: savedAddresses[0].AddressStr,
                    city: savedAddresses[0].City,
                    district: savedAddresses[0].District,
                    zipCode: savedAddresses[0].ZipCode
                }));
            } else {
                setSelectedAddressIndex(-1);
                setSelectedBillingAddressIndex(-1);
            }
        });
    }, []);

    const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const { name, value, type, checked } = e.target;

        if (name === "tcNo") {
            const onlyNums = value.replace(/[^0-9]/g, '');
            if (onlyNums.length <= 11) {
                setFormData(prev => ({ ...prev, [name]: onlyNums }));
            }
            return;
        }

        setFormData(prev => ({ ...prev, [name]: type === 'checkbox' ? checked : value }));
    };

    const handleAddressSelectionChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const index = parseInt(e.target.value);
        setSelectedAddressIndex(index);

        if (index >= 0 && savedAddresses[index]) {
            const addr = savedAddresses[index];
            setFormData(prev => ({
                ...prev,
                address: addr.AddressStr,
                city: addr.City,
                district: addr.District,
                zipCode: addr.ZipCode
            }));
            setSaveNewAddress(false);
        } else {
            setFormData(prev => ({
                ...prev,
                address: "",
                city: "",
                district: "",
                zipCode: ""
            }));
        }
    };

    const handleBillingAddressSelectionChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const index = parseInt(e.target.value);
        setSelectedBillingAddressIndex(index);

        if (index >= 0 && savedAddresses[index]) {
            const addr = savedAddresses[index];
            setFormData(prev => ({
                ...prev,
                billingAddress: addr.AddressStr,
                billingCity: addr.City,
                billingDistrict: addr.District,
                billingZipCode: addr.ZipCode
            }));
        } else {
            setFormData(prev => ({
                ...prev,
                billingAddress: "",
                billingCity: "",
                billingDistrict: "",
                billingZipCode: ""
            }));
        }
    };

    const handleUseSameAddressChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const checked = e.target.checked;
        setUseSameAddressForBilling(checked);

        if (!checked && selectedBillingAddressIndex >= 0 && savedAddresses[selectedBillingAddressIndex]) {
            const addr = savedAddresses[selectedBillingAddressIndex];
            setFormData(prev => ({
                ...prev,
                billingAddress: addr.AddressStr,
                billingCity: addr.City,
                billingDistrict: addr.District,
                billingZipCode: addr.ZipCode
            }));
        }
    };

    const handleNext = () => {
        setActiveStep((prev) => prev + 1);
    };

    const handleBack = () => {
        setActiveStep((prev) => prev - 1);
    };

    const handleConfirmOrder = async () => {
        const requiredFields = [
            'tcNo',
            'phone',
            'address',
            'city',
            'district',
            'zipCode',
            ...(!useSameAddressForBilling ? [
                'billingAddress',
                'billingCity',
                'billingDistrict',
                'billingZipCode'
            ] : []),
            'cardName',
            'cardNumber',
            'expireMonth',
            'expireYear',
            'cvv'
        ];

        const isFormValid = requiredFields.every(field => formData[field as keyof typeof formData].toString().trim() !== "");

        if (!isFormValid) {
            setNotification({ open: true, message: "Lütfen tüm zorunlu alanları doldurunuz.", severity: "error" });
            return;
        }

        if (formData.tcNo.length !== 11) {
            setNotification({ open: true, message: "Lütfen 11 haneli geçerli bir TC Kimlik No giriniz.", severity: "error" });
            return;
        }

        const phoneRegex = /^(?:\+90|0)?5\d{2}[ \-]?\d{3}[ \-]?\d{2}[ \-]?\d{2}$/;
        if (!phoneRegex.test(formData.phone)) {
            setNotification({ open: true, message: "Lütfen geçerli bir telefon numarası giriniz.", severity: "error" });
            return;
        }

        setLoading(true);

        if (selectedAddressIndex === -1 && saveNewAddress) {
            await userService.addAddress({
                AddressName: "Kayıtlı Profil Adresim",
                AddressStr: formData.address,
                City: formData.city,
                District: formData.district,
                Country: "Türkiye",
                ZipCode: formData.zipCode
            }).catch(console.error);
        }

        // API çağrısı 
        const orderData: CreateOrderRequest = {
            Tckn: formData.tcNo,
            PhoneNumber: formData.phone,
            CardData: {
                CardName: formData.cardName,
                CardNumber: formData.cardNumber,
                ExpireYear: formData.expireYear,
                ExpireMonth: formData.expireMonth,
                Cvc: formData.cvv
            },
            BillingAddress: {
                AddressStr: useSameAddressForBilling ? formData.address : formData.billingAddress,
                City: useSameAddressForBilling ? formData.city : formData.billingCity,
                Country: "Türkiye",
                ZipCode: useSameAddressForBilling ? formData.zipCode : formData.billingZipCode
            },
            ShippingAddress: {
                AddressStr: formData.address,
                City: formData.city,
                Country: "Türkiye",
                ZipCode: formData.zipCode
            },
            SecurePay: formData.use3DSecure,
            SaveUserAddress: saveNewAddress
        };

        const result = await orderService.createOrder(orderData);

        if (!result) {
            setNotification({ open: true, message: "Siparişiniz oluşturulamadı.", severity: "error" });
            setLoading(false);
            return;
        }

        setNotification({ open: true, message: "Siparişiniz başarıyla oluşturuldu.", severity: "success" });
        clearCart();
        setTimeout(() => {
            navigate('/profile/orders');
        }, 2000);
    };

    const handleCloseNotification = () => {
        setNotification({ ...notification, open: false });
    };

    return (
        <Container maxWidth="xl" sx={{ py: { xs: 4, md: 8 }, bgcolor: 'background.default' }}>
            <Box sx={{ mb: 6, textAlign: 'center' }}>
                <Typography variant="h3" fontWeight="900" gutterBottom>
                    Güvenli Ödeme
                </Typography>
                <Typography variant="subtitle1" color="text.secondary">
                    Siparişinizi tamamlamak için lütfen bilgilerinizi eksiksiz giriniz.
                </Typography>
            </Box>

            <Grid container spacing={4}>
                {/* Sol Taraf: Form Alanı */}
                <Grid size={{ xs: 12, md: 8 }}>
                    <Paper elevation={0} sx={{ p: { xs: 3, md: 5 }, borderRadius: 4, border: '1px solid', borderColor: 'divider' }}>

                        {/* Stepper */}
                        <Stepper activeStep={activeStep} alternativeLabel sx={{ mb: 4 }}>
                            {steps.map((label) => (
                                <Step key={label}>
                                    <StepLabel>{label}</StepLabel>
                                </Step>
                            ))}
                        </Stepper>

                        <Divider sx={{ mb: 4 }} />

                        {/* Step 0: Kimlik Bilgileri */}
                        {activeStep === 0 && (
                            <Box>
                                <Typography variant="h6" fontWeight="bold" sx={{ mb: 3 }}>Kimlik Bilgileri</Typography>
                                <Grid container spacing={2}>
                                    <Grid size={{ xs: 12, md: 6 }}>
                                        <TextField
                                            fullWidth size="small" label="TC Kimlik No" name="tcNo"
                                            value={formData.tcNo} onChange={handleInputChange}
                                            variant="outlined" required
                                            error={formData.tcNo.length > 0 && formData.tcNo.length !== 11}
                                            helperText={formData.tcNo.length > 0 && formData.tcNo.length !== 11 ? "11 haneli olmalıdır" : ""}
                                            slotProps={{ htmlInput: { maxLength: 11 } }}

                                        />
                                    </Grid>
                                    <Grid size={{ xs: 12, md: 6 }}>
                                        <TextField
                                            fullWidth size="small" label="Telefon Numarası" name="phone"
                                            value={formData.phone} onChange={handleInputChange}
                                            placeholder="5XX XXX XX XX"
                                            variant="outlined" required
                                        />
                                    </Grid>
                                </Grid>
                            </Box>
                        )}

                        {/* Step 1: Teslimat Adresi */}
                        {activeStep === 1 && (
                            <Box>
                                <Typography variant="h6" fontWeight="bold" sx={{ mb: 3 }}>Teslimat Adresi</Typography>

                                {savedAddresses.length === 0 && (
                                    <Alert severity="info" sx={{ mb: 3 }}>
                                        Kayıtlı teslimat adresiniz bulunmuyor. Lütfen aşağıdan yeni adres bilgilerinizi giriniz.
                                    </Alert>
                                )}

                                {savedAddresses.length > 0 && (
                                    <FormControl component="fieldset" sx={{ mb: 3, width: '100%' }}>
                                        <RadioGroup
                                            value={selectedAddressIndex.toString()}
                                            onChange={handleAddressSelectionChange}
                                        >
                                            {savedAddresses.map((addr, idx) => (
                                                <Paper key={idx} variant="outlined" sx={{ mb: 1, p: 1, borderRadius: 2, bgcolor: selectedAddressIndex === idx ? 'action.hover' : 'transparent' }}>
                                                    <FormControlLabel
                                                        value={idx.toString()}
                                                        control={<Radio />}
                                                        label={
                                                            <Box sx={{ ml: 1, py: 1 }}>
                                                                <Typography fontWeight="bold">{addr.AddressName || `Adres ${idx + 1}`}</Typography>
                                                                <Typography variant="body2" color="text.secondary">{addr.AddressStr}</Typography>
                                                                <Typography variant="body2" color="text.secondary">{addr.District} / {addr.City}</Typography>
                                                            </Box>
                                                        }
                                                        sx={{ width: '100%', m: 0 }}
                                                    />
                                                </Paper>
                                            ))}
                                            <Paper variant="outlined" sx={{ p: 1, borderRadius: 2, bgcolor: selectedAddressIndex === -1 ? 'action.hover' : 'transparent' }}>
                                                <FormControlLabel
                                                    value="-1"
                                                    control={<Radio />}
                                                    label={<Typography fontWeight="bold" sx={{ ml: 1 }}>Yeni / Farklı Bir Adres Gir</Typography>}
                                                    sx={{ width: '100%', m: 0 }}
                                                />
                                            </Paper>
                                        </RadioGroup>
                                    </FormControl>
                                )}

                                {(selectedAddressIndex === -1 || savedAddresses.length === 0) && (
                                    <Grid container spacing={2}>
                                        <Grid size={{ xs: 12 }}>
                                            <TextField
                                                fullWidth size="small" label="Açık Adres" name="address"
                                                value={formData.address} onChange={handleInputChange}
                                                multiline rows={2} required

                                            />
                                        </Grid>
                                        <Grid size={{ xs: 12, sm: 4 }}>
                                            <TextField fullWidth size="small" label="İl" name="city" value={formData.city} onChange={handleInputChange} required />
                                        </Grid>
                                        <Grid size={{ xs: 12, sm: 4 }}>
                                            <TextField fullWidth size="small" label="İlçe" name="district" value={formData.district} onChange={handleInputChange} required />
                                        </Grid>
                                        <Grid size={{ xs: 12, sm: 4 }}>
                                            <TextField fullWidth size="small" label="Posta Kodu" name="zipCode" value={formData.zipCode} onChange={handleInputChange} required />
                                        </Grid>
                                        <Grid size={{ xs: 12 }}>
                                            <FormControlLabel
                                                control={<Checkbox checked={saveNewAddress} onChange={(e) => setSaveNewAddress(e.target.checked)} color="primary" />}
                                                label="Bu adresi profilime kaydet"
                                            />
                                        </Grid>
                                    </Grid>
                                )}
                            </Box>
                        )}

                        {/* Step 2: Fatura Adresi */}
                        {activeStep === 2 && (
                            <Box>
                                <Typography variant="h6" fontWeight="bold" sx={{ mb: 3 }}>Fatura Adresi</Typography>

                                <FormControlLabel
                                    control={<Checkbox checked={useSameAddressForBilling} onChange={handleUseSameAddressChange} color="primary" />}
                                    label={<Typography fontWeight="medium">Fatura adresim teslimat adresim ile aynı olsun</Typography>}
                                    sx={{ mb: 3 }}
                                />

                                {!useSameAddressForBilling && (
                                    <>
                                        {savedAddresses.length === 0 && (
                                            <Alert severity="info" sx={{ mb: 3 }}>
                                                Kayıtlı adresiniz bulunmuyor. Lütfen aşağıdan fatura adresinizi giriniz.
                                            </Alert>
                                        )}

                                        {savedAddresses.length > 0 && (
                                            <FormControl component="fieldset" sx={{ mb: 3, width: '100%' }}>
                                                <RadioGroup
                                                    value={selectedBillingAddressIndex.toString()}
                                                    onChange={handleBillingAddressSelectionChange}
                                                >
                                                    {savedAddresses.map((addr, idx) => (
                                                        <Paper key={idx} variant="outlined" sx={{ mb: 1, p: 1, borderRadius: 2, bgcolor: selectedBillingAddressIndex === idx ? 'action.hover' : 'transparent' }}>
                                                            <FormControlLabel
                                                                value={idx.toString()}
                                                                control={<Radio />}
                                                                label={
                                                                    <Box sx={{ ml: 1, py: 1 }}>
                                                                        <Typography fontWeight="bold">{addr.AddressName || `Adres ${idx + 1}`}</Typography>
                                                                        <Typography variant="body2" color="text.secondary">{addr.AddressStr}</Typography>
                                                                        <Typography variant="body2" color="text.secondary">{addr.District} / {addr.City}</Typography>
                                                                    </Box>
                                                                }
                                                                sx={{ width: '100%', m: 0 }}
                                                            />
                                                        </Paper>
                                                    ))}
                                                    <Paper variant="outlined" sx={{ p: 1, borderRadius: 2, bgcolor: selectedBillingAddressIndex === -1 ? 'action.hover' : 'transparent' }}>
                                                        <FormControlLabel
                                                            value="-1"
                                                            control={<Radio />}
                                                            label={<Typography fontWeight="bold" sx={{ ml: 1 }}>Yeni / Farklı Bir Fatura Adresi Gir</Typography>}
                                                            sx={{ width: '100%', m: 0 }}
                                                        />
                                                    </Paper>
                                                </RadioGroup>
                                            </FormControl>
                                        )}

                                        {(selectedBillingAddressIndex === -1 || savedAddresses.length === 0) && (
                                            <Grid container spacing={2}>
                                                <Grid size={{ xs: 12 }}>
                                                    <TextField
                                                        fullWidth size="small" label="Fatura Adresi" name="billingAddress"
                                                        value={formData.billingAddress} onChange={handleInputChange}
                                                        multiline rows={2} required

                                                    />
                                                </Grid>
                                                <Grid size={{ xs: 12, sm: 4 }}>
                                                    <TextField fullWidth size="small" label="İl" name="billingCity" value={formData.billingCity} onChange={handleInputChange} required />
                                                </Grid>
                                                <Grid size={{ xs: 12, sm: 4 }}>
                                                    <TextField fullWidth size="small" label="İlçe" name="billingDistrict" value={formData.billingDistrict} onChange={handleInputChange} required />
                                                </Grid>
                                                <Grid size={{ xs: 12, sm: 4 }}>
                                                    <TextField fullWidth size="small" label="Posta Kodu" name="billingZipCode" value={formData.billingZipCode} onChange={handleInputChange} required />
                                                </Grid>
                                            </Grid>
                                        )}
                                    </>
                                )}
                            </Box>
                        )}

                        {/* Step 3: Ödeme Bilgileri */}
                        {activeStep === 3 && (
                            <Box>
                                <Typography variant="h6" fontWeight="bold" sx={{ mb: 3 }}>Kart Bilgileri</Typography>
                                <Grid container spacing={2}>
                                    <Grid size={{ xs: 12 }}>
                                        <TextField fullWidth size="small" label="Kart Üzerindeki İsim" name="cardName" value={formData.cardName} onChange={handleInputChange} required />
                                    </Grid>
                                    <Grid size={{ xs: 12 }}>
                                        <TextField fullWidth size="small" label="Kart Numarası" name="cardNumber" value={formData.cardNumber} onChange={handleInputChange} placeholder="0000 0000 0000 0000" required />
                                    </Grid>
                                    <Grid size={{ xs: 12, sm: 4 }}>
                                        <TextField fullWidth size="small" label="Ay" name="expireMonth" value={formData.expireMonth} onChange={handleInputChange} placeholder="AA" required />
                                    </Grid>
                                    <Grid size={{ xs: 12, sm: 4 }}>
                                        <TextField fullWidth size="small" label="Yıl" name="expireYear" value={formData.expireYear} onChange={handleInputChange} placeholder="YYYY" required />
                                    </Grid>
                                    <Grid size={{ xs: 12, sm: 4 }}>
                                        <TextField fullWidth size="small" label="CVV" name="cvv" value={formData.cvv} onChange={handleInputChange} slotProps={{ htmlInput: { maxLength: 3 } }} required />
                                    </Grid>
                                    <Grid size={{ xs: 12 }}>
                                        <FormControlLabel
                                            control={<Checkbox checked={formData.use3DSecure} onChange={handleInputChange} name="use3DSecure" color="primary" />}
                                            label={
                                                <Box display="flex" alignItems="center">
                                                    <SecurityIcon color="success" sx={{ mr: 1, fontSize: 20 }} />
                                                    <Typography variant="body2">3D Secure ile Güvenli Ödeme Yapmak İstiyorum</Typography>
                                                </Box>
                                            }
                                        />
                                    </Grid>
                                </Grid>
                            </Box>
                        )}

                        {/* Navigation Buttons */}
                        <Box sx={{ display: 'flex', justifyContent: 'space-between', mt: 5, pt: 2, borderTop: '1px solid', borderColor: 'divider' }}>
                            <Button
                                color="inherit"
                                disabled={activeStep === 0}
                                onClick={handleBack}
                                sx={{ px: 3 }}
                            >
                                Geri
                            </Button>

                            {activeStep === steps.length - 1 ? (
                                <Typography variant="body2" color="success.main" fontWeight="bold" sx={{ alignSelf: 'center' }}>
                                    Bilgiler tamam. Sağdaki panelden siparişi onaylayabilirsiniz.
                                </Typography>
                            ) : (
                                <Button
                                    variant="contained"
                                    onClick={handleNext}
                                    sx={{ px: 4 }}
                                >
                                    Devam Et
                                </Button>
                            )}
                        </Box>
                    </Paper>
                </Grid>

                {/* Sağ Taraf: Özet ve Onay */}
                <Grid size={{ xs: 12, md: 4 }}>
                    <Paper elevation={0} sx={{ p: 4, border: '1px solid', borderColor: 'divider', position: 'sticky', top: 100, bgcolor: 'grey.50' }}>
                        <Typography variant="h6" fontWeight="bold" gutterBottom>
                            Sipariş Özeti
                        </Typography>
                        <Divider sx={{ my: 3 }} />
                        <Box display="flex" justifyContent="space-between" mb={2}>
                            <Typography color="text.secondary">Ara Toplam</Typography>
                            <Typography fontWeight="bold">{totalPrice.toLocaleString('tr-TR')} ₺</Typography>
                        </Box>
                        <Box display="flex" justifyContent="space-between" mb={2}>
                            <Typography color="text.secondary">Kargo Ücreti</Typography>
                            <Typography fontWeight="bold" color="success.main">Ücretsiz</Typography>
                        </Box>
                        <Divider sx={{ my: 3 }} />
                        <Box display="flex" justifyContent="space-between" alignItems="center" sx={{ mb: 4 }}>
                            <Typography variant="h6" fontWeight="bold">
                                Ödenecek Tutar
                            </Typography>
                            <Typography variant="h4" color="primary" fontWeight="bold">
                                {totalPrice.toLocaleString('tr-TR')} ₺
                            </Typography>
                        </Box>

                        <Button
                            variant="contained"
                            size="large"
                            fullWidth
                            color="secondary"
                            disabled={loading || activeStep !== steps.length - 1}
                            onClick={handleConfirmOrder}
                            sx={{ 
                                py: 2, 
                                fontSize: '1.1rem', 
                                fontWeight: 'bold', 
                                textTransform: 'none', 
                                mb: 2,
                                boxShadow: '0 4px 14px 0 rgba(249, 115, 22, 0.39)',
                                '&:hover': {
                                    boxShadow: '0 6px 20px 0 rgba(249, 115, 22, 0.45)',
                                }
                            }}
                        >
                            {loading ? 'İşleniyor...' : 'Siparişi Onayla'}
                        </Button>

                        {activeStep !== steps.length - 1 && (
                            <Typography variant="caption" color="error.main" display="block" textAlign="center" sx={{ mb: 2 }}>
                                * Onaylamak için tüm adımları tamamlayın.
                            </Typography>
                        )}

                        <Button
                            variant="outlined"
                            size="large"
                            fullWidth
                            onClick={() => navigate('/cart')}
                            sx={{ py: 1.5, textTransform: 'none', borderWidth: 2, '&:hover': { borderWidth: 2 } }}
                        >
                            Sepete Dön
                        </Button>

                        <Box sx={{ mt: 4, textAlign: 'center', p: 2, bgcolor: 'white', borderRadius: 2, border: '1px dashed', borderColor: 'success.main' }}>
                            <SecurityIcon color="success" sx={{ fontSize: 32, mb: 1 }} />
                            <Typography variant="body2" color="text.secondary" fontWeight="medium">
                                Ödemeleriniz 256-bit SSL şifreleme ile güvence altındadır.
                            </Typography>
                        </Box>
                    </Paper>
                </Grid>
            </Grid>

            <Snackbar open={notification.open} autoHideDuration={6000} onClose={handleCloseNotification} anchorOrigin={{ vertical: 'bottom', horizontal: 'center' }}>
                <Alert onClose={handleCloseNotification} severity={notification.severity} sx={{ width: '100%' }}>
                    {notification.message}
                </Alert>
            </Snackbar>
        </Container>
    );
}

export default OrderSummary;
