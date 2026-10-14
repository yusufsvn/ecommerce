import React, { useEffect, useState } from 'react';
import {
  Box,
  Typography,
  Paper,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Chip,
  MenuItem,
  Select,
  FormControl,
  IconButton,
  Collapse,
  Avatar,
  CircularProgress,
  Grid,
  Backdrop,
  Pagination,
  TextField,
  Button,
  InputLabel,
  Checkbox
} from '@mui/material';
import type { SelectChangeEvent } from '@mui/material';
import KeyboardArrowDownIcon from '@mui/icons-material/KeyboardArrowDown';
import KeyboardArrowUpIcon from '@mui/icons-material/KeyboardArrowUp';
import { orderService } from '../../Services/orderService';
import type { OrderForAdmin, OrderItem, OrderProducts } from '../../models/ordermodel';
import { statusOptions, formatStatusToTurkish, mapTurkishToEnglishStatus, getStatusColor } from '../../utils/statusMapper';

// Alt satır bileşeni (Detaylar)
const Row = ({ row, handleStatusChange, isSelected, handleSelect }: { row: OrderForAdmin, handleStatusChange: (id: string, newStatus: string) => void, isSelected: boolean, handleSelect: (id: string) => void }) => {
  const [open, setOpen] = useState(false);
  const [orderItems, setOrderItems] = useState<OrderItem | null>(null);
  const [isLoading, setIsLoading] = useState(false);

  const handleToggle = () => {
    const willOpen = !open;
    setOpen(willOpen);

    // Fetch details only when opening the row and if we haven't fetched them yet
    if (willOpen && !orderItems) {
      setIsLoading(true);
      orderService.getOrderDetails(row.order_id).then((res) => {
        if (res.success && res.orderItems) {
          setOrderItems(res.orderItems);
        }
        setIsLoading(false);
      });
    }
  };

  const onChange = (e: SelectChangeEvent) => {
    handleStatusChange(row.order_id, e.target.value);
  };

  return (
    <React.Fragment>
      <Backdrop
        sx={{ color: '#fff', zIndex: (theme) => theme.zIndex.drawer + 1 }}
        open={isLoading}
      >
        <CircularProgress color="inherit" />
      </Backdrop>
      <TableRow sx={{ '& > *': { borderBottom: 'unset' } }} selected={isSelected}>
        <TableCell padding="checkbox">
          <Checkbox
            color="primary"
            checked={isSelected}
            onChange={() => handleSelect(row.order_id)}
          />
        </TableCell>
        <TableCell>
          <IconButton
            aria-label="expand row"
            size="small"
            onClick={handleToggle}
          >
            {open ? <KeyboardArrowUpIcon /> : <KeyboardArrowDownIcon />}
          </IconButton>
        </TableCell>
        <TableCell component="th" scope="row">
          {row.order_id}
        </TableCell>
        <TableCell>{`${row.name} ${row.surname}`}</TableCell>
        <TableCell>{new Date(row.created_at).toLocaleDateString('tr-TR')}</TableCell>
        <TableCell>{row.total_amount} TL</TableCell>
        <TableCell>
          <Chip label={row.status} color={getStatusColor(row.status) as any} size="small" />
        </TableCell>
        <TableCell>
          <FormControl size="small" sx={{ minWidth: 140 }}>
            <Select
              value={row.status}
              onChange={onChange}
              displayEmpty
              inputProps={{ 'aria-label': 'Without label' }}
              sx={{
                bgcolor: 'background.paper',
                borderRadius: 1,
                '& .MuiSelect-select': { py: 1 },
              }}
            >
              {statusOptions.map((option) => (
                <MenuItem key={option} value={option}>
                  {option}
                </MenuItem>
              ))}
            </Select>
          </FormControl>
        </TableCell>
      </TableRow>
      <TableRow>
        <TableCell style={{ paddingBottom: 0, paddingTop: 0 }} colSpan={8}>
          <Collapse in={open} timeout="auto" unmountOnExit>
            <Box sx={{ margin: 2, p: 2, bgcolor: 'background.default', borderRadius: 2 }}>
              <Typography variant="h6" gutterBottom component="div">
                Sipariş Detayları
              </Typography>
              {isLoading ? (
                <Box sx={{ display: 'flex', justifyContent: 'center', py: 3 }} />
              ) : !orderItems || orderItems.products.length === 0 ? (
                <Typography variant="body2" color="text.secondary" align="center" sx={{ py: 3 }}>
                  Ürün detayı bulunamadı.
                </Typography>
              ) : (
                <Grid container spacing={3} alignItems="stretch">
                  <Grid size={{ xs: 12, md: 8 }}>
                    <Table size="small" aria-label="purchases" sx={{ height: '100%' }}>
                      <TableHead>
                        <TableRow>
                          <TableCell align="center">Görsel</TableCell>
                          <TableCell>Ürün Adı</TableCell>
                          <TableCell align="center">Miktar</TableCell>
                          <TableCell align="right">Fiyat</TableCell>
                        </TableRow>
                      </TableHead>
                      <TableBody>
                        {orderItems.products.map((itemRow: OrderProducts, index: number) => (
                          <TableRow key={index}>
                            <TableCell align="center">
                              <Avatar variant="rounded" src={itemRow.thumbnail_image} alt={itemRow.product_name} sx={{ width: 48, height: 48, margin: 'auto' }} />
                            </TableCell>
                            <TableCell component="th" scope="row">
                              {itemRow.product_name}
                            </TableCell>
                            <TableCell align="center">{itemRow.quantity}</TableCell>
                            <TableCell align="right">{itemRow.price_at_purchase} TL</TableCell>
                          </TableRow>
                        ))}
                      </TableBody>
                    </Table>
                  </Grid>
                  <Grid size={{ xs: 12, md: 4 }}>
                    <Box sx={{ p: 2, bgcolor: 'background.paper', borderRadius: 1, border: '1px solid #e0e0e0', height: '100%', display: 'flex', flexDirection: 'column', justifyContent: 'center' }}>
                      <Box>
                        <Typography variant="subtitle2" color="text.primary" gutterBottom sx={{ fontWeight: 'bold' }}>
                          Teslimat Adresi
                        </Typography>
                        <Typography variant="body2" color="text.secondary" paragraph sx={{ whiteSpace: 'pre-wrap' }}>
                          {orderItems.shipping_address || 'Belirtilmemiş'}
                        </Typography>
                      </Box>

                      <Box sx={{ mt: 1 }}>
                        <Typography variant="subtitle2" color="text.primary" gutterBottom sx={{ fontWeight: 'bold' }}>
                          Fatura Adresi
                        </Typography>
                        <Typography variant="body2" color="text.secondary" sx={{ whiteSpace: 'pre-wrap' }}>
                          {orderItems.billing_address || 'Belirtilmemiş'}
                        </Typography>
                      </Box>
                    </Box>
                  </Grid>
                </Grid>
              )}
            </Box>
          </Collapse>
        </TableCell>
      </TableRow>
    </React.Fragment>
  );
};

const AdminOrders = () => {
  const [orders, setOrders] = useState<OrderForAdmin[]>([]);
  const [totalCount, setTotalCount] = useState<number>(0);
  const [page, setPage] = useState<number>(1);
  const limit = 5;

  const [statusFilter, setStatusFilter] = useState<string>('');
  const [searchTerm, setSearchTerm] = useState<string>('');
  const [appliedFilters, setAppliedFilters] = useState({ status: '', search: '' });

  const [selectedOrders, setSelectedOrders] = useState<string[]>([]);
  const [bulkStatus, setBulkStatus] = useState<string>('');

  const handleFilter = () => {
    setPage(1);
    setAppliedFilters({ status: statusFilter, search: searchTerm });
  };

  const handleSelectAllClick = (event: React.ChangeEvent<HTMLInputElement>) => {
    if (event.target.checked) {
      const newSelected = orders.map((n) => n.order_id);
      setSelectedOrders(newSelected);
      return;
    }
    setSelectedOrders([]);
  };

  const handleSelect = (id: string) => {
    const selectedIndex = selectedOrders.indexOf(id);
    let newSelected: string[] = [];

    if (selectedIndex === -1) {
      newSelected = newSelected.concat(selectedOrders, id);
    } else if (selectedIndex === 0) {
      newSelected = newSelected.concat(selectedOrders.slice(1));
    } else if (selectedIndex === selectedOrders.length - 1) {
      newSelected = newSelected.concat(selectedOrders.slice(0, -1));
    } else if (selectedIndex > 0) {
      newSelected = newSelected.concat(
        selectedOrders.slice(0, selectedIndex),
        selectedOrders.slice(selectedIndex + 1),
      );
    }
    setSelectedOrders(newSelected);
  };

  const handleBulkStatusChange = async () => {
    if (!bulkStatus || selectedOrders.length === 0) return;

    const englishState = mapTurkishToEnglishStatus(bulkStatus);

    const response = await orderService.updateMultiStatus(selectedOrders, englishState);
    if (response) {
      setOrders((prevOrders) =>
        prevOrders.map((order) =>
          selectedOrders.includes(order.order_id) ? { ...order, status: bulkStatus } : order
        )
      );
      setSelectedOrders([]);
      setBulkStatus('');
    }
  };

  const handleStatusChange = async (id: string, newStatus: string) => {

    const state = mapTurkishToEnglishStatus(newStatus);
    const success = await orderService.updateOrderStatus(id, state);
    if (success) {
      setOrders((prevOrders) =>
        prevOrders.map((order) =>
          order.order_id === id ? { ...order, status: newStatus } : order
        )
      );
    } else {
      console.error("Sipariş durumu güncellenirken bir hata oluştu.");
    }
  };

  useEffect(() => {
    const fetchOrders = async () => {
      const offset = (page - 1) * limit;
      let res;

      try {
        if (appliedFilters.status) {
          res = await orderService.getOrdersByStatus(limit, offset, mapTurkishToEnglishStatus(appliedFilters.status));
        } else {
          res = await orderService.getAllOrders(offset, limit);
        }

        let mappedOrders = (res.data || []).map(order => ({
          ...order,
          status: formatStatusToTurkish(order.status)
        }));

        if (appliedFilters.search) {
          const lowerSearch = appliedFilters.search.toLowerCase();
          mappedOrders = mappedOrders.filter((order: OrderForAdmin) =>
            (order.order_id && order.order_id.toLowerCase().includes(lowerSearch)) ||
            (order.name && order.surname && `${order.name} ${order.surname}`.toLowerCase().includes(lowerSearch))
          );
        }

        setOrders(mappedOrders);
        setTotalCount(res.count || 0);
        setSelectedOrders([]);
      } catch (error) {
        console.error("Siparişler yüklenirken hata oluştu:", error);
      }
    };

    fetchOrders();
  }, [page, appliedFilters]);

  return (
    <Box sx={{ width: '100%' }}>
      <Typography variant="h4" component="h2" gutterBottom sx={{ fontWeight: 600, color: 'text.primary', mb: 3 }}>
        Sipariş Yönetimi
      </Typography>

      <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 3 }}>
        <Box sx={{ display: 'flex', gap: 2, alignItems: 'center' }}>
          <TextField
            label="Sipariş No veya İsim Ara"
            variant="outlined"
            size="small"
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
            sx={{ minWidth: 250 }}
          />
          <FormControl size="small" sx={{ minWidth: 200 }}>
            <InputLabel>Durum Filtresi</InputLabel>
            <Select
              value={statusFilter}
              label="Durum Filtresi"
              onChange={(e) => setStatusFilter(e.target.value)}
            >
              <MenuItem value="">Tümü</MenuItem>
              {statusOptions.map((option) => (
                <MenuItem key={option} value={option}>
                  {option}
                </MenuItem>
              ))}
            </Select>
          </FormControl>
          <Button variant="contained" onClick={handleFilter} sx={{ height: 40 }}>
            Filtrele
          </Button>
        </Box>

        {selectedOrders.length > 0 && (
          <Box sx={{ display: 'flex', gap: 2, alignItems: 'center', bgcolor: 'primary.light', p: 1, borderRadius: 1 }}>
            <Typography variant="body2" sx={{ color: 'white', fontWeight: 'bold', ml: 1 }}>
              {selectedOrders.length} Seçili
            </Typography>
            <FormControl size="small" sx={{ minWidth: 150, bgcolor: 'background.paper', borderRadius: 1 }}>
              <Select
                value={bulkStatus}
                onChange={(e) => setBulkStatus(e.target.value)}
                displayEmpty
              >
                <MenuItem value="" disabled>Durum Seçin</MenuItem>
                {statusOptions.map((option) => (
                  <MenuItem key={option} value={option}>{option}</MenuItem>
                ))}
              </Select>
            </FormControl>
            <Button
              variant="contained"
              color="secondary"
              onClick={handleBulkStatusChange}
              disabled={!bulkStatus}
              sx={{ height: 40 }}
            >
              Uygula
            </Button>
          </Box>
        )}
      </Box>

      <TableContainer component={Paper} elevation={3} sx={{ overflow: 'hidden' }}>
        <Table aria-label="collapsible table">
          <TableHead sx={{ bgcolor: 'grey.50' }}>
            <TableRow>
              <TableCell padding="checkbox">
                <Checkbox
                  color="primary"
                  indeterminate={selectedOrders.length > 0 && selectedOrders.length < orders.length}
                  checked={orders.length > 0 && selectedOrders.length === orders.length}
                  onChange={handleSelectAllClick}
                />
              </TableCell>
              <TableCell />
              <TableCell>Sipariş Kodu</TableCell>
              <TableCell>Müşteri</TableCell>
              <TableCell>Tarih</TableCell>
              <TableCell>Toplam Tutar</TableCell>
              <TableCell>Durum</TableCell>
              <TableCell>Durum Güncelle</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {orders.map((order) => {
              const isSelected = selectedOrders.indexOf(order.order_id) !== -1;
              return (
                <Row
                  key={order.order_id}
                  row={order}
                  handleStatusChange={handleStatusChange}
                  isSelected={isSelected}
                  handleSelect={handleSelect}
                />
              );
            })}
          </TableBody>
        </Table>
      </TableContainer>

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
    </Box>
  );
};

export default AdminOrders;
