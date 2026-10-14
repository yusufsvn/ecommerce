export interface Registeruser {
    name: string
    second_name: string
    mail: string
    password: string
    password_check: string
}

export interface LoginAdmin {
    mail: string;
    password: string;
}

export interface LoginUserRes {
    name: string;
    surname: string;
}
