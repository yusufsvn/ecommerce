export interface addressRequest {
    AddressId?: string;
    AddressName: string;
    AddressStr: string;
    City: string;
    District: string;
    Country: string;
    ZipCode: string;
}

export interface AddressUserRes {
    Addresses: addressRequest[];
}
