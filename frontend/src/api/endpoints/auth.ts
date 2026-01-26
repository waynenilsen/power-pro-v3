import { post, get } from '../client';

// Request types

export interface RegisterRequest {
  email: string;
  password: string;
  name?: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

// Response types

export interface UserResponse {
  id: string;
  email: string;
  name: string;
  createdAt: string;
  updatedAt: string;
}

export interface LoginResponse {
  token: string;
  expiresAt: string;
  user: UserResponse;
}

// Auth endpoints

export async function register(request: RegisterRequest): Promise<UserResponse> {
  return post<UserResponse, RegisterRequest>('/auth/register', request);
}

export async function login(request: LoginRequest): Promise<LoginResponse> {
  return post<LoginResponse, LoginRequest>('/auth/login', request);
}

export async function logout(): Promise<void> {
  return post<void>('/auth/logout');
}

export async function me(): Promise<UserResponse> {
  return get<UserResponse>('/auth/me');
}
