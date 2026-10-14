import { Box, Card, CardContent, CardMedia, Typography, IconButton } from '@mui/material'
import ShoppingCartOutlinedIcon from '@mui/icons-material/ShoppingCartOutlined';
import CheckIcon from '@mui/icons-material/Check';
import { useNavigate } from 'react-router';
import { useRef, useState, useEffect } from 'react';

import type { GetProducts, productprops } from '../models/ProductResponse';

function ProductCard(productprop: productprops) {
  const navigate = useNavigate();
  const cardRef = useRef<HTMLDivElement>(null);
  const [detectedDark, setDetectedDark] = useState<boolean>(false);

  const product: GetProducts = productprop.product;
  const cardImagePath = product.thumbnail_image;

  // İçinde bulunduğu ana kutunun (container) arka plan rengini otomatik algıla
  useEffect(() => {
    if (productprop.isDarkContainer !== undefined) return;
    if (!cardRef.current) return;

    let el: HTMLElement | null = cardRef.current.parentElement;
    while (el && el !== document.body) {
      const bg = window.getComputedStyle(el).backgroundColor;
      if (bg && bg !== 'transparent' && bg !== 'rgba(0, 0, 0, 0)') {
        const match = bg.match(/\d+/g);
        if (match && match.length >= 3) {
          const [r, g, b] = match.map(Number);
          // YIQ parlaklık formülü (0: en koyu, 255: en açık)
          const brightness = (r * 299 + g * 587 + b * 114) / 1000;
          setDetectedDark(brightness < 128);
          break;
        }
      }
      el = el.parentElement;
    }
  }, [productprop.isDarkContainer]);

  // Ana kutu koyu ise açık kart modu, ana kutu açık ise koyu (main renk) kart modu
  const inDarkContainer = productprop.isDarkContainer !== undefined
    ? productprop.isDarkContainer
    : detectedDark;

  // Renk seçimleri:
  // 1) Ana kutu KOYU ise -> Yazı kısmı AÇIK (Beyaz), Yazılar KOYU (#1E3A8A)
  // 2) Ana kutu AÇIK ise -> Yazı kısmı KOYU (Main: #1E3A8A), Yazılar AÇIK (Beyaz & Highlight)
  const contentBgColor = inDarkContainer ? 'background.paper' : 'primary.main';
  const titleColor = inDarkContainer ? 'primary.main' : '#FFFFFF';
  const priceColor = inDarkContainer ? 'secondary.main' : '#38BDF8';
  const borderColor = inDarkContainer ? 'rgba(255, 255, 255, 0.2)' : 'divider';

  return (
    <Card
      ref={cardRef}
      elevation={0}
      sx={{
        height: '100%',
        display: 'flex',
        flexDirection: 'column',
        cursor: 'pointer',
        border: '1px solid',
        borderColor: borderColor,
        bgcolor: 'background.paper',
        borderRadius: 3,
        overflow: 'hidden',
        transition: 'all 0.3s cubic-bezier(0.4, 0, 0.2, 1)',
        '&:hover': {
          transform: 'translateY(-6px)',
          boxShadow: inDarkContainer
            ? '0 16px 30px rgba(0, 0, 0, 0.35)'
            : '0 16px 30px rgba(30, 58, 138, 0.2)',
        },
      }}
      onClick={() => navigate(`/product/${product.product_id}`)}
    >
      {/* Görsel Kutusu: Her zaman temiz beyaz zemin */}
      <Box sx={{ p: 2, bgcolor: '#FFFFFF', display: 'flex', justifyContent: 'center', alignItems: 'center' }}>
        <CardMedia
          component="img"
          height="210"
          image={cardImagePath}
          alt={product.product_name}
          sx={{ objectFit: 'contain', transition: 'transform 0.3s', '&:hover': { transform: 'scale(1.05)' } }}
        />
      </Box>

      {/* Yazıların bulunduğu kısım: Arka plana göre dinamik renk */}
      <CardContent
        sx={{
          flexGrow: 1,
          display: 'flex',
          flexDirection: 'column',
          justifyContent: 'space-between',
          p: 2.2,
          bgcolor: contentBgColor,
          transition: 'background-color 0.3s ease',
        }}
      >
        <Box>
          <Typography
            gutterBottom
            variant="subtitle1"
            component="h2"
            noWrap
            sx={{
              fontWeight: 700,
              color: titleColor,
              fontSize: '1rem',
            }}
          >
            {product.product_name}
          </Typography>
        </Box>

        <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-end', mt: 2 }}>
          <Typography
            variant="h6"
            sx={{
              fontWeight: 800,
              color: priceColor,
              letterSpacing: '-0.5px',
            }}
          >
            {product.price.toFixed(2)} ₺
          </Typography>

          <IconButton
            disabled={productprop.isInCart}
            onClick={(e) => {
              e.stopPropagation();
              return productprop.addtocart({
                product_id: product.product_id,
                product_name: product.product_name,
                price: product.price,
                stock: product.stock,
                thumbnail_image: product.thumbnail_image,
                quantity: 1
              });
            }}
            sx={{
              bgcolor: productprop.isInCart ? 'success.main' : 'secondary.main',
              color: 'white',
              borderRadius: '50%',
              width: 40,
              height: 40,
              boxShadow: '0 2px 8px rgba(249, 115, 22, 0.35)',
              '&:hover': {
                bgcolor: productprop.isInCart ? 'success.dark' : 'secondary.dark',
              },
              '&.Mui-disabled': {
                bgcolor: 'success.light',
                color: 'white',
              }
            }}
          >
            {productprop.isInCart ? <CheckIcon fontSize="small" /> : <ShoppingCartOutlinedIcon fontSize="small" />}
          </IconButton>
        </Box>
      </CardContent>
    </Card>
  );
}

export default ProductCard;


