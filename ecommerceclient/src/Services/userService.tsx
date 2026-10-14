import { api } from "../Interceptors/httpinterceptor";
import type { addressRequest, AddressUserRes } from "../models/addressmodel";
import type { Registeruser, LoginAdmin, LoginUserRes } from "../models/userModels";

const publicUrl = "https://ecommerceapi.yusufsevinc.work/api/user";
const adminUrl = "https://ecommerceapi.yusufsevinc.work/api/admin";
export const userService = {


    registerUser: async (data: Registeruser) => {
        const response = await api.post(publicUrl + "/registeruser", data, {
            headers: {
                'Content-Type': 'application/json'
            }
        })
        return response.status

    },

    loginAdmin: async (data: LoginAdmin) => {
        const response = await api.post(adminUrl + "/login", data, {
            headers: {
                'Content-Type': 'application/json'
            },
            withCredentials: true
        })
        return response.status
    }
    ,
    loginUser: async (data: LoginAdmin) => {
        const response = await api.post<LoginUserRes>(publicUrl + "/login", data, {
            headers: {
                'Content-Type': 'application/json'
            },
            withCredentials: true
        })
        return response.data
    },
    loginWithGoogle: async (token: string) => {
        const response = await api.post<LoginUserRes>(publicUrl + "/loginwithgoogle", { 
            returnIdToken: token }, {
            headers: {
                'Content-Type': 'application/json'
            }
        })
        return response.data
    },

    getAllUserAddress: async () => {
        const response = await api.get<AddressUserRes>(publicUrl + "/getaddresses", {
            headers: {
                'Content-Type': 'application/json'
            },
            withCredentials: true
        })
        return response.data
    },
    updateAddress: async (data: addressRequest) => {
        const response = await api.put(publicUrl + "/updateaddress", data, {
            headers: {
                'Content-Type': 'application/json'
            },
            withCredentials: true
        })
        return response.status
    },

    deleteAddress: async (id: string) => {
        const response = await api.delete(publicUrl + "/deleteaddress", {
            params: {
                addressId: id
            },
            headers: {
                'Content-Type': 'application/json'
            },
            withCredentials: true
        })
        return response.status
    },

    addAddress: async (data: addressRequest) => {
        const response = await api.post(publicUrl + "/createaddress", data, {
            headers: {
                'Content-Type': 'application/json'
            },
            withCredentials: true
        })
        return response.status
    }
}
// client de api istekleri yerleştirilsin