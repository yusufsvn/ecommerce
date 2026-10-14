import type { CreateOrderRequest, OrderItem, OrderResponse, OrderResult } from "../models/ordermodel";
import { api } from "../Interceptors/httpinterceptor";

const baseUrl = "https://ecommerceapi.yusufsevinc.work/api";

export const orderService = {
    createOrder: async (orderData: CreateOrderRequest) => {
        try {
            const response = await api.post(baseUrl + "/user/createorder", orderData, {
                headers: {
                    'Content-Type': 'application/json'
                },
                withCredentials: true
            });
            if (response.status === 200 || response.status === 201) {
                return true;
            }
            return false;
        } catch (error) {
            console.error("Error creating order:", error);
            return false;
        }
    },

    getUserOrders: async (): Promise<{ orders: OrderResponse[], success: boolean }> => {

        const response = await api.get<OrderResponse[]>(baseUrl + "/user/getorderbyuserid", {
            headers: {
                'Content-Type': 'application/json'
            },
            withCredentials: true
        })

        if (response.status === 200) {
            return { orders: response.data, success: true }
        }
        else {
            return { orders: [], success: false };
        }
    },

    getOrderDetails: async (orderId: string): Promise<{ orderItems: OrderItem | undefined, success: boolean }> => {
        const response = await api.get<OrderItem>(baseUrl + "/user/getorderdetails?orderId=" + orderId, {
            headers: {
                'Content-Type': 'application/json'
            },
            withCredentials: true
        })

        if (response.status === 200) {
            return { orderItems: response.data, success: true }
        }
        else {
            return { orderItems: undefined, success: false };
        }
    },

    getAllOrders: async (offset: number, limit: number): Promise<OrderResult> => {
        const params: any = {
            limit: limit,
            offset: offset
        };

        const response = await api.get<OrderResult>(baseUrl + "/admin/getallorders", {
            headers: {
                'Content-Type': 'application/json'
            },
            params: params,
            withCredentials: true
        })

        if (response.status === 200) {
            return response.data
        }
        else {
            return { count: 0, data: [] };
        }
    },

    updateOrderStatus: async (orderId: string, state: string) => {
        try {
            const response = await api.put(baseUrl + "/admin/updateorderstate", {
                order_id: orderId,
                new_state: state
            }, {
                headers: {
                    'Content-Type': 'application/json'
                },
                withCredentials: true
            });
            if (response.status === 200 || response.status === 201) {
                return true;
            }
            return false;
        } catch (error) {
            console.error("Error updating order status:", error);
            return false;
        }
    },

    updateMultiStatus: async (orderIds: string[], newState: string) => {
        try {
            const response = await api.put(baseUrl + "/admin/updatestates", {
                OrderIds: orderIds,
                New_state: newState
            }, {
                headers: {
                    'Content-Type': 'application/json'
                },
                withCredentials: true
            });
            if (response.status === 200 || response.status === 201) {
                return true;
            }
            return false;
        } catch (error) {
            console.error("Error updating multi status:", error);
            return false;
        }
    },

    getOrdersByStatus: async (limit: number, offset: number, status: string) => {
        const params: any = {
            limit: limit,
            offset: offset,
            status: status
        };
        const response = await api.get<OrderResult>(baseUrl + "/admin/getordersbystatus", {
            headers: {
                'Content-Type': 'application/json'
            },
            params: params,
            withCredentials: true
        })

        if (response.status === 200) {
            return response.data
        }
        else {
            return { count: 0, data: [] };
        }
    },

}


