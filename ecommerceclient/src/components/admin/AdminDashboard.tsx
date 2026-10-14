import { Box, Container, Typography, Paper, Stack, Avatar, Grid, Button, Chip, CircularProgress } from '@mui/material';
import PeopleAltIcon from '@mui/icons-material/PeopleAlt';
import InventoryIcon from '@mui/icons-material/Inventory';
import ShoppingCartIcon from '@mui/icons-material/ShoppingCart';
import AttachMoneyIcon from '@mui/icons-material/AttachMoney';
import ArrowForwardIcon from '@mui/icons-material/ArrowForward';
import ReceiptLongOutlinedIcon from '@mui/icons-material/ReceiptLongOutlined';
import { useNavigate } from 'react-router';
import { useEffect, useState } from 'react';
import { api } from '../../Interceptors/httpinterceptor';
import { orderService } from '../../Services/orderService';
import type { OrderForAdmin } from '../../models/ordermodel';
import { formatStatusToTurkish, getStatusColor } from '../../utils/statusMapper';

function AdminDashboard() {
    const navigate = useNavigate();
    const [userCount, setUserCount] = useState<number | string>('...');
    const [productCount, setProductCount] = useState<number | string>('...');
    const [waitingOrderCount, setWaitingOrderCount] = useState<number | string>('...');
    const [recentOrders, setRecentOrders] = useState<OrderForAdmin[]>([]);
    const [ordersLoading, setOrdersLoading] = useState<boolean>(true);

    useEffect(() => {
        const fetchDashboardData = async () => {
            try {
                const userRes = await api.get('admin/getusercount', { withCredentials: true });
                setUserCount(userRes.data);

                const productRes = await api.get('admin/getproductcount', { withCredentials: true });
                setProductCount(productRes.data);

                const orderRes = await api.get('admin/getordercountforwait', { withCredentials: true });
                setWaitingOrderCount(orderRes.data);
            } catch (error) {
                console.error("Dashboard verileri alınırken hata oluştu", error);
            }
        };

        const fetchRecentOrders = async () => {
            try {
                setOrdersLoading(true);
                const orderResult = await orderService.getAllOrders(0, 5);
                if (orderResult && orderResult.data) {
                    setRecentOrders(orderResult.data);
                }
            } catch (error) {
                console.error("Son siparişler alınırken hata oluştu", error);
            } finally {
                setOrdersLoading(false);
            }
        };

        fetchDashboardData();
        fetchRecentOrders();
    }, []);

    const formatDate = (dateStr: string) => {
        if (!dateStr) return '-';
        try {
            const d = new Date(dateStr);
            return isNaN(d.getTime()) ? dateStr : d.toLocaleDateString('tr-TR', {
                day: '2-digit',
                month: 'short',
                year: 'numeric',
                hour: '2-digit',
                minute: '2-digit'
            });
        } catch {
            return dateStr;
        }
    };

    const formatPrice = (amount: string | number) => {
        const num = Number(amount);
        if (isNaN(num)) return `${amount} ₺`;
        return `${num.toLocaleString('tr-TR', { minimumFractionDigits: 2, maximumFractionDigits: 2 })} ₺`;
    };

    const stats = [
        { title: 'Toplam Kullanıcı', value: userCount, icon: <PeopleAltIcon fontSize="large" />, color: 'success.main', bgColor: 'success.light' },
        { title: 'Toplam Ürün', value: productCount, icon: <InventoryIcon fontSize="large" />, color: 'info.main', bgColor: 'info.light', path: '/admin/products' },
        { title: 'Bekleyen Sipariş', value: waitingOrderCount, icon: <ShoppingCartIcon fontSize="large" />, color: 'warning.main', bgColor: 'warning.light', path: '/admin/orders' },
        { title: 'Aylık Gelir', value: '₺45,200', icon: <AttachMoneyIcon fontSize="large" />, color: 'secondary.main', bgColor: 'secondary.light' },
    ];

    return (
        <Container maxWidth="xl" sx={{ mt: 2, mb: 4 }}>
            <Box sx={{ mb: 4 }}>
                <Typography variant="h4" fontWeight="bold" color="text.primary">
                    Dashboard
                </Typography>
                <Typography variant="subtitle1" color="text.secondary">
                    Admin paneline hoş geldiniz. Mağazanızın güncel durumunu buradan takip edebilirsiniz.
                </Typography>
            </Box>

            {/* İstatistik Kartları */}
            <Grid container spacing={3} sx={{ mb: 4 }}>
                {stats.map((stat, index) => (
                    <Grid size={{ xs: 12, sm: 6, md: 3 }} key={index}>
                        <Paper
                            onClick={() => stat.path && navigate(stat.path)}
                            sx={{
                                p: 3,
                                display: 'flex',
                                alignItems: 'center',
                                gap: 2,
                                transition: 'transform 0.2s',
                                cursor: stat.path ? 'pointer' : 'default',
                                '&:hover': {
                                    transform: stat.path ? 'translateY(-4px)' : 'none',
                                    boxShadow: stat.path ? 4 : 1
                                }
                            }}
                        >
                            <Avatar sx={{ bgcolor: stat.bgColor, color: stat.color, width: 60, height: 60 }}>
                                {stat.icon}
                            </Avatar>
                            <Box>
                                <Typography variant="h5" fontWeight="bold">
                                    {stat.value}
                                </Typography>
                                <Typography variant="body2" color="text.secondary" fontWeight="medium">
                                    {stat.title}
                                </Typography>
                            </Box>
                        </Paper>
                    </Grid>
                ))}
            </Grid>

            <Grid container spacing={3}>
                {/* Son Siparişler */}
                <Grid size={{ xs: 12, md: 8 }}>
                    <Paper sx={{ p: 3, height: '100%', display: 'flex', flexDirection: 'column' }}>
                        <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2.5 }}>
                            <Box>
                                <Typography variant="h6" fontWeight="bold">
                                    Son Siparişler
                                </Typography>
                                <Typography variant="body2" color="text.secondary">
                                    Mağazanızdan verilen en son siparişlerin güncel durumu
                                </Typography>
                            </Box>
                            <Button
                                size="small"
                                endIcon={<ArrowForwardIcon fontSize="small" />}
                                onClick={() => navigate('/admin/orders')}
                                sx={{ textTransform: 'none', fontWeight: 600 }}
                            >
                                Tümünü İncele
                            </Button>
                        </Box>

                        {ordersLoading ? (
                            <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', py: 8, flexGrow: 1 }}>
                                <CircularProgress size={36} />
                            </Box>
                        ) : recentOrders.length === 0 ? (
                            <Box sx={{ textAlign: 'center', py: 8, flexGrow: 1, display: 'flex', flexDirection: 'column', alignItems: 'center', justifyContent: 'center' }}>
                                <ReceiptLongOutlinedIcon sx={{ fontSize: 48, color: 'text.disabled', mb: 1 }} />
                                <Typography variant="body2" color="text.secondary">
                                    Henüz kayıtlı bir sipariş bulunmuyor.
                                </Typography>
                            </Box>
                        ) : (
                            <Stack spacing={1.5} sx={{ flexGrow: 1 }}>
                                {recentOrders.map((order) => {
                                    const trStatus = formatStatusToTurkish(order.status);
                                    const statusColor = getStatusColor(trStatus);
                                    const initials = `${order.name?.charAt(0) || ''}${order.surname?.charAt(0) || ''}`.toUpperCase() || 'SP';

                                    return (
                                        <Box
                                            key={order.order_id}
                                            onClick={() => navigate('/admin/orders')}
                                            sx={{
                                                display: 'flex',
                                                alignItems: 'center',
                                                justifyContent: 'space-between',
                                                p: 2,
                                                bgcolor: 'background.default',
                                                borderRadius: 2,
                                                border: '1px solid',
                                                borderColor: 'divider',
                                                cursor: 'pointer',
                                                transition: 'all 0.2s ease',
                                                '&:hover': {
                                                    bgcolor: 'action.hover',
                                                    borderColor: 'primary.light',
                                                    transform: 'translateY(-2px)',
                                                    boxShadow: '0 4px 12px rgba(0, 0, 0, 0.05)',
                                                },
                                            }}
                                        >
                                            <Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
                                                <Avatar
                                                    sx={{
                                                        width: 44,
                                                        height: 44,
                                                        bgcolor: 'primary.light',
                                                        color: 'primary.contrastText',
                                                        fontSize: '0.875rem',
                                                        fontWeight: 700,
                                                    }}
                                                >
                                                    {initials}
                                                </Avatar>
                                                <Box>
                                                    <Typography variant="subtitle2" fontWeight="bold" color="text.primary">
                                                        Sipariş #{order.order_id.slice(0, 8)}
                                                    </Typography>
                                                    <Typography variant="body2" color="text.secondary">
                                                        {order.name} {order.surname}
                                                    </Typography>
                                                    <Typography variant="caption" color="text.disabled">
                                                        {formatDate(order.created_at)}
                                                    </Typography>
                                                </Box>
                                            </Box>

                                            <Box sx={{ textAlign: 'right', display: 'flex', flexDirection: 'column', alignItems: 'flex-end', gap: 0.75 }}>
                                                <Typography variant="subtitle1" color="primary.main" fontWeight="bold">
                                                    {formatPrice(order.total_amount)}
                                                </Typography>
                                                <Chip
                                                    label={trStatus}
                                                    color={statusColor as any}
                                                    size="small"
                                                    sx={{ fontWeight: 600, fontSize: '0.75rem', height: 24 }}
                                                />
                                            </Box>
                                        </Box>
                                    );
                                })}
                            </Stack>
                        )}
                    </Paper>
                </Grid>

                {/* Hızlı İşlemler */}
                <Grid size={{ xs: 12, md: 4 }}>
                    <Paper sx={{ p: 3, height: '100%' }}>
                        <Typography variant="h6" fontWeight="bold" sx={{ mb: 2 }}>
                            Hızlı İşlemler
                        </Typography>
                        <Typography variant="body2" color="text.secondary" sx={{ mb: 3 }}>
                            Sık kullanılan admin görevlerine hızlıca erişin.
                        </Typography>
                        <Stack spacing={2}>
                            <Paper
                                variant="outlined"
                                onClick={() => navigate('/admin/create-product')}
                                sx={{ p: 2, textAlign: 'center', cursor: 'pointer', '&:hover': { bgcolor: 'primary.50', borderColor: 'primary.main' } }}
                            >
                                <Typography color="primary.main" fontWeight="medium">+ Yeni Ürün Ekle</Typography>
                            </Paper>
                            <Paper
                                variant="outlined"
                                onClick={() => navigate('/admin/categories')}
                                sx={{ p: 2, textAlign: 'center', cursor: 'pointer', '&:hover': { bgcolor: 'primary.50', borderColor: 'primary.main' } }}
                            >
                                <Typography color="primary.main" fontWeight="medium">Kategorileri Yönet</Typography>
                            </Paper>
                            <Paper
                                variant="outlined"
                                onClick={() => navigate('/admin/orders')}
                                sx={{ p: 2, textAlign: 'center', cursor: 'pointer', '&:hover': { bgcolor: 'primary.50', borderColor: 'primary.main' } }}
                            >
                                <Typography color="primary.main" fontWeight="medium">Siparişleri İncele</Typography>
                            </Paper>
                        </Stack>
                    </Paper>
                </Grid>
            </Grid>
        </Container>
    )
}

export default AdminDashboard
