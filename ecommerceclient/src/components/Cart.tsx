import { Avatar, Box, Button, Container, Divider, IconButton, Stack, Typography, Paper, Grid } from "@mui/material"
import type { cartItem, cartProps } from "../models/cartmodel"
import { Add, DeleteOutline, Remove, ShoppingBagOutlined } from "@mui/icons-material"
import { useNavigate } from "react-router";

function Cart(cartprops: cartProps) {
  const navigate = useNavigate();
  const cartItems: cartItem[] = cartprops.products
  const totalPrice: number = cartprops.calculateTotalPrice()

  return (
    <Container maxWidth="xl" sx={{ py: { xs: 4, md: 8 } }}>
      <Typography variant="h4" component="h1" fontWeight="bold" gutterBottom sx={{ mb: 4, display: 'flex', alignItems: 'center', gap: 1 }}>
        <ShoppingBagOutlined fontSize="large" color="primary" /> Sepetim ({cartItems.length} Ürün)
      </Typography>

      {cartItems.length === 0 ? (
        <Box sx={{ py: 10, textAlign: 'center', bgcolor: 'background.paper', borderRadius: 2, border: '1px solid', borderColor: 'divider' }}>
          <ShoppingBagOutlined sx={{ fontSize: 80, color: 'grey.300', mb: 2 }} />
          <Typography variant="h5" color="text.secondary" fontWeight="medium" gutterBottom>
            Sepetiniz şu an boş
          </Typography>
          <Typography color="text.secondary" sx={{ mb: 4 }}>
            Sepetinize ürün eklemek için hemen alışverişe başlayın.
          </Typography>
          <Button variant="contained" size="large" onClick={() => navigate('/products')} sx={{ px: 4, py: 1.5 }}>
            Alışverişe Başla
          </Button>
        </Box>
      ) : (
        <Grid container spacing={4}>
          {/* Ürün Listesi */}
          <Grid size={{ xs: 12, md: 8 }}>
            <Stack spacing={3}>
              {cartItems.map((item) => (
                <Paper key={item.product_id} elevation={0} sx={{ p: 3, border: '1px solid', borderColor: 'divider', transition: 'box-shadow 0.3s', '&:hover': { boxShadow: '0 4px 20px rgba(0,0,0,0.05)' } }}>
                  <Grid container spacing={3} alignItems="center">
                    {/* Görsel */}
                    <Grid size={{ xs: 12, sm: 2 }} sx={{ display: 'flex', justifyContent: 'center' }}>
                      <Avatar
                        src={item.thumbnail_image}
                        alt={item.product_name}
                        variant="rounded"
                        sx={{ width: 100, height: 100, bgcolor: 'grey.50' }}
                      />
                    </Grid>

                    {/* Bilgiler */}
                    <Grid size={{ xs: 12, sm: 5 }}>
                      <Typography variant="subtitle2" color="success.main" fontWeight="bold" sx={{ mb: 0.5 }}>
                        Stokta Var
                      </Typography>
                      <Typography variant="h6" fontWeight="bold" sx={{ mb: 1, lineHeight: 1.2 }}>
                        {item.product_name}
                      </Typography>
                      <Typography variant="h6" color="primary" fontWeight="bold">
                        {item.price.toLocaleString('tr-TR')} ₺
                      </Typography>
                    </Grid>

                    {/* Miktar */}
                    <Grid size={{ xs: 6, sm: 3 }} sx={{ display: 'flex', justifyContent: { xs: 'flex-start', sm: 'center' } }}>
                      <Box sx={{ display: 'flex', alignItems: 'center', border: '1px solid', borderColor: 'divider', borderRadius: 2, p: 0.5 }}>
                        <IconButton size="small" onClick={() => cartprops.handleDecrement(item.product_id)} sx={{ bgcolor: 'grey.50', '&:hover': { bgcolor: 'grey.200' } }}>
                          <Remove fontSize="small" />
                        </IconButton>
                        <Typography sx={{ px: 2, minWidth: 30, textAlign: 'center', fontWeight: 'bold' }}>
                          {item.quantity}
                        </Typography>
                        <IconButton size="small" onClick={() => cartprops.handleIncrement(item.product_id)} sx={{ bgcolor: 'grey.50', '&:hover': { bgcolor: 'grey.200' } }}>
                          <Add fontSize="small" />
                        </IconButton>
                      </Box>
                    </Grid>

                    {/* Sil Butonu */}
                    <Grid size={{ xs: 6, sm: 2 }} sx={{ display: 'flex', justifyContent: 'flex-end' }}>
                      <IconButton color="error" onClick={() => cartprops.deleteProduct(item.product_id)} sx={{ bgcolor: 'error.50', '&:hover': { bgcolor: 'error.100' } }}>
                        <DeleteOutline />
                      </IconButton>
                    </Grid>
                  </Grid>
                </Paper>
              ))}
            </Stack>
          </Grid>

          {/* Sipariş Özeti */}
          <Grid size={{ xs: 12, md: 4 }}>
            <Box sx={{ position: { md: 'sticky' }, top: 100 }}>
              <Paper elevation={0} sx={{ p: 4, border: '1px solid', borderColor: 'divider', bgcolor: 'grey.50' }}>
                <Typography variant="h6" fontWeight="bold" gutterBottom>
                  Sipariş Özeti
                </Typography>
                <Divider sx={{ my: 3 }} />

                <Stack spacing={2} sx={{ mb: 3 }}>
                  <Box display="flex" justifyContent="space-between" alignItems="center">
                    <Typography color="text.secondary">Ara Toplam</Typography>
                    <Typography fontWeight="bold">{totalPrice.toLocaleString('tr-TR')} ₺</Typography>
                  </Box>
                  <Box display="flex" justifyContent="space-between" alignItems="center">
                    <Typography color="text.secondary">Kargo</Typography>
                    <Typography fontWeight="bold" color="success.main">Ücretsiz</Typography>
                  </Box>
                </Stack>

                <Divider sx={{ my: 3 }} />

                <Box display="flex" justifyContent="space-between" alignItems="center" sx={{ mb: 4 }}>
                  <Typography variant="h6" fontWeight="bold" color="text.primary">
                    Genel Toplam
                  </Typography>
                  <Typography variant="h4" color="primary.main" fontWeight="bold">
                    {totalPrice.toLocaleString('tr-TR')} ₺
                  </Typography>
                </Box>

                <Button
                  variant="contained"
                  size="large"
                  fullWidth
                  color="secondary"
                  sx={{ 
                    py: 2, 
                    fontWeight: 'bold', 
                    fontSize: '1.1rem', 
                    textTransform: 'none',
                    boxShadow: '0 4px 14px 0 rgba(249, 115, 22, 0.39)',
                    '&:hover': {
                      boxShadow: '0 6px 20px 0 rgba(249, 115, 22, 0.45)',
                    }
                  }}
                  onClick={() => navigate("/order-summary")}
                >
                  Alışverişi Tamamla
                </Button>

                <Box sx={{ mt: 3, textAlign: 'center' }}>
                  <Typography variant="caption" color="text.secondary">
                    256-bit SSL Sertifikası ile Güvenli Alışveriş
                  </Typography>
                </Box>
              </Paper>
            </Box>
          </Grid>
        </Grid>
      )}
    </Container>
  )
}

export default Cart
