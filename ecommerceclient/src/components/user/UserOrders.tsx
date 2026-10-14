import { useEffect, useState } from 'react';
import {
    Container, Typography, Box, Chip, Button, Collapse, Avatar, CircularProgress, Paper, Divider, Stack, Grid
} from '@mui/material';
import { useNavigate } from 'react-router';
import KeyboardArrowDownIcon from '@mui/icons-material/KeyboardArrowDown';
import KeyboardArrowUpIcon from '@mui/icons-material/KeyboardArrowUp';
import ArrowBackIcon from '@mui/icons-material/ArrowBack';
import ReceiptLongIcon from '@mui/icons-material/ReceiptLong';
import { orderService } from '../../Services/orderService';
import type { OrderItem, OrderResponse } from '../../models/ordermodel';
import { formatStatusToTurkish, getStatusColor } from '../../utils/statusMapper';

function OrderCard({ order }: { order: OrderResponse }) {
    const [open, setOpen] = useState(false);
    const [orderItem, setOrderItem] = useState<OrderItem | undefined>(undefined);
    const [isLoading, setIsLoading] = useState(false);

    const handleToggle = () => {
        const willOpen = !open;
        setOpen(willOpen);

        // Fetch details only when opening the row and if we haven't fetched them yet
        if (willOpen && !orderItem) {
            setIsLoading(true);
            orderService.getOrderDetails(order.order_code).then((res) => {
                if (res.success) {
                    setOrderItem(res.orderItems);
                }
                setIsLoading(false);
            });
        }
    };


    return (
        <Paper elevation={0} sx={{ mb: 3, border: '1px solid', borderColor: 'divider', overflow: 'hidden' }}>
            {/* Header / Summary Section */}
            <Box sx={{ p: 3, bgcolor: open ? 'grey.50' : 'white', transition: 'background-color 0.3s' }}>
                <Grid container spacing={3} alignItems="center">
                    <Grid size={{ xs: 6, sm: 2.5 }}>
                        <Typography variant="caption" color="text.secondary" display="block">
                            Sipariş Kodu
                        </Typography>
                        <Typography variant="body1" fontWeight="bold">
                            {order.order_code}
                        </Typography>
                    </Grid>
                    <Grid size={{ xs: 6, sm: 2 }}>
                        <Typography variant="caption" color="text.secondary" display="block">
                            Sipariş Tarihi
                        </Typography>
                        <Typography variant="body1" fontWeight="bold">
                            {new Date(order.created_at).toLocaleDateString('tr-TR')}
                        </Typography>
                    </Grid>
                    <Grid size={{ xs: 6, sm: 2 }}>
                        <Typography variant="caption" color="text.secondary" display="block">
                            Sipariş Özeti
                        </Typography>
                        <Typography variant="body1" fontWeight="bold">
                            {order.total_quantity} Ürün
                        </Typography>
                    </Grid>
                    <Grid size={{ xs: 6, sm: 2.5 }}>
                        <Typography variant="caption" color="text.secondary" display="block">
                            Toplam Tutar
                        </Typography>
                        <Typography variant="body1" color="primary" fontWeight="bold">
                            {order.total_amount} TL
                        </Typography>
                    </Grid>
                    <Grid size={{ xs: 12, sm: 3 }} sx={{ display: 'flex', justifyContent: { xs: 'flex-start', sm: 'flex-end' }, alignItems: 'center', gap: 2 }}>
                        <Chip
                            label={formatStatusToTurkish(order.status)}
                            color={getStatusColor(order.status) as any}
                            size="small"
                            sx={{ fontWeight: 'bold' }}
                        />
                        <Button
                            size="small"
                            variant="outlined"
                            onClick={handleToggle}
                            endIcon={open ? <KeyboardArrowUpIcon /> : <KeyboardArrowDownIcon />}
                            sx={{ textTransform: 'none' }}
                        >
                            Detaylar
                        </Button>
                    </Grid>
                </Grid>
            </Box>

            {/* Details Section */}
            <Collapse in={open} timeout="auto" unmountOnExit>
                <Divider />
                <Box sx={{ p: { xs: 2, md: 4 }, bgcolor: 'white' }}>
                    {isLoading ? (
                        <Box sx={{ display: 'flex', justifyContent: 'center', py: 5 }}>
                            <CircularProgress />
                        </Box>
                    ) : !orderItem ? (
                        <Typography variant="body1" color="text.secondary" align="center" sx={{ py: 5 }}>
                            Ürün detayı bulunamadı.
                        </Typography>
                    ) : (
                        <Grid container spacing={4}>
                            <Grid size={{ xs: 12, md: 8 }}>
                                <Typography variant="h6" fontWeight="bold" gutterBottom sx={{ mb: 3 }}>
                                    Teslim Edilecek Ürünler
                                </Typography>
                                <Stack spacing={3}>
                                    {orderItem.products.map(product =>
                                        <Box key={product.product_id} sx={{ display: 'flex', alignItems: 'center', p: 2, border: '1px solid', borderColor: 'divider' }}>
                                            <Avatar variant="rounded" src={product.thumbnail_image} alt={product.product_name} sx={{ width: 80, height: 80, mr: 3, bgcolor: 'grey.50' }} />
                                            <Box sx={{ flexGrow: 1 }}>
                                                <Typography variant="subtitle1" fontWeight="bold">
                                                    {product.product_name}
                                                </Typography>
                                                <Typography variant="body2" color="text.secondary">
                                                    Miktar: {product.quantity}
                                                </Typography>
                                            </Box>
                                            <Typography variant="h6" color="primary" fontWeight="bold">
                                                {product.price_at_purchase} TL
                                            </Typography>
                                        </Box>
                                    )
                                    }
                                </Stack>
                            </Grid>

                            <Grid size={{ xs: 12, md: 4 }}>
                                <Typography variant="h6" fontWeight="bold" gutterBottom sx={{ mb: 3 }}>
                                    Teslimat ve Fatura
                                </Typography>
                                <Paper elevation={0} sx={{ p: 3, bgcolor: 'grey.50', height: '100%' }}>
                                    <Box sx={{ mb: 3 }}>
                                        <Typography variant="subtitle2" color="text.primary" fontWeight="bold" gutterBottom>
                                            Teslimat Adresi
                                        </Typography>
                                        <Typography variant="body2" color="text.secondary" sx={{ whiteSpace: 'pre-wrap', lineHeight: 1.6 }}>
                                            {orderItem.shipping_address || 'Belirtilmemiş'}
                                        </Typography>
                                    </Box>
                                    <Divider sx={{ my: 2 }} />
                                    <Box>
                                        <Typography variant="subtitle2" color="text.primary" fontWeight="bold" gutterBottom>
                                            Fatura Adresi
                                        </Typography>
                                        <Typography variant="body2" color="text.secondary" sx={{ whiteSpace: 'pre-wrap', lineHeight: 1.6 }}>
                                            {orderItem.billing_address || 'Belirtilmemiş'}
                                        </Typography>
                                    </Box>
                                </Paper>
                            </Grid>
                        </Grid>
                    )}
                </Box>
            </Collapse>
        </Paper>
    );
}

export default function UserOrders() {
    const navigate = useNavigate();
    const [orders, setOrders] = useState<OrderResponse[]>([]);

    useEffect(() => {
        // get api call
        orderService.getUserOrders().then((res) => {
            if (res.success) {
                setOrders(res.orders);
            }
        });
    }, []);

    return (
        <Container component="main" maxWidth="lg" sx={{ py: { xs: 4, md: 8 } }}>
            <Box sx={{ mb: 4 }}>
                <Button
                    startIcon={<ArrowBackIcon />}
                    onClick={() => navigate('/profile')}
                    sx={{ color: 'text.secondary', '&:hover': { bgcolor: 'grey.100' }, px: 2 }}
                >
                    Profile Dön
                </Button>
            </Box>

            <Box sx={{ display: 'flex', alignItems: 'center', mb: 6, justifyContent: 'center', flexDirection: 'column' }}>
                <Box sx={{ p: 2, borderRadius: '50%', bgcolor: 'primary.50', mb: 2 }}>
                    <ReceiptLongIcon color="primary" sx={{ fontSize: 40 }} />
                </Box>
                <Typography component="h1" variant="h3" fontWeight="900" align="center" gutterBottom>
                    Siparişlerim
                </Typography>
                <Typography variant="subtitle1" color="text.secondary" align="center">
                    Tüm geçmiş ve devam eden siparişlerinizi buradan detaylıca inceleyebilirsiniz.
                </Typography>
            </Box>

            {orders.length === 0 ? (
                <Paper elevation={0} sx={{ p: 8, textAlign: 'center', border: '1px dashed', borderColor: 'divider' }}>
                    <Typography variant="h6" color="text.secondary" gutterBottom>
                        Henüz hiç siparişiniz bulunmamaktadır.
                    </Typography>
                    <Typography color="text.secondary" sx={{ mb: 4 }}>
                        Size özel ürünleri keşfetmek için mağazamıza göz atın.
                    </Typography>
                    <Button variant="contained" size="large" onClick={() => navigate('/products')} sx={{ px: 4, py: 1.5 }}>
                        Alışverişe Başla
                    </Button>
                </Paper>
            ) : (
                <Box>
                    {orders.map((order) => (
                        <OrderCard key={order.order_id} order={order} />
                    ))}
                </Box>
            )}
        </Container>
    );
}
