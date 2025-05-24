import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, BehaviorSubject, tap } from 'rxjs';
import { Router } from '@angular/router';

export interface User {
  id: string;
  email: string;
  name: string;
  picture: string;
  role: string;
  created_at: string;
  updated_at: string;
}

export interface LoginResponse {
  token: string;
  refresh_token: string;
  user: User;
  expires_in: number;
}

export interface RegistrationRequest {
  email: string;
  name: string;
  password: string;
}

@Injectable({
  providedIn: 'root'
})
export class AuthService {
  private apiUrl = 'http://localhost:8080/api/auth';
  private currentUserSubject = new BehaviorSubject<User | null>(null);
  public currentUser$ = this.currentUserSubject.asObservable();
  
  constructor(
    private http: HttpClient,
    private router: Router
  ) {
    // Check if user is already logged in
    const token = localStorage.getItem('access_token');
    if (token) {
      this.getProfile().subscribe({
        next: (user) => this.currentUserSubject.next(user),
        error: () => this.logout()
      });
    }
  }

  // Register a new user
  register(email: string, name: string, password: string): Observable<LoginResponse> {
    return this.http.post<LoginResponse>(`${this.apiUrl}/register`, {
      email,
      name,
      password
    }).pipe(
      tap(response => {
        this.storeTokens(response);
        this.currentUserSubject.next(response.user);
      })
    );
  }

  // Initiate Google OAuth login
  googleLogin(): Observable<{ auth_url: string }> {
    return this.http.get<{ auth_url: string }>(`${this.apiUrl}/google/login`);
  }

  // Initiate Google OAuth registration
  googleRegister(): Observable<{ auth_url: string }> {
    return this.http.get<{ auth_url: string }>(`${this.apiUrl}/google/register`);
  }

  // Handle OAuth callback (for both login and registration)
  handleCallback(code: string, state: string): Observable<LoginResponse> {
    return this.http.get<LoginResponse>(`${this.apiUrl}/google/callback`, {
      params: { code, state }
    }).pipe(
      tap(response => {
        this.storeTokens(response);
        this.currentUserSubject.next(response.user);
      })
    );
  }

  // Refresh token
  refreshToken(): Observable<LoginResponse> {
    const refreshToken = localStorage.getItem('refresh_token');
    if (!refreshToken) {
      throw new Error('No refresh token available');
    }

    return this.http.post<LoginResponse>(`${this.apiUrl}/refresh`, {
      refresh_token: refreshToken
    }).pipe(
      tap(response => {
        this.storeTokens(response);
        this.currentUserSubject.next(response.user);
      })
    );
  }

  // Get user profile
  getProfile(): Observable<User> {
    return this.http.get<User>(`${this.apiUrl}/profile`);
  }

  // Logout
  logout(): void {
    this.http.post(`${this.apiUrl}/logout`, {}).subscribe({
      complete: () => {
        localStorage.removeItem('access_token');
        localStorage.removeItem('refresh_token');
        localStorage.removeItem('token_expires_at');
        this.currentUserSubject.next(null);
        this.router.navigate(['/login']);
      }
    });
  }

  // Check if user is authenticated
  isAuthenticated(): boolean {
    const token = localStorage.getItem('access_token');
    const expiresAt = localStorage.getItem('token_expires_at');
    
    if (!token || !expiresAt) {
      return false;
    }

    const now = new Date().getTime();
    const expiry = parseInt(expiresAt);
    
    return now < expiry;
  }

  // Get access token
  getAccessToken(): string | null {
    return localStorage.getItem('access_token');
  }

  // Get current user
  getCurrentUser(): User | null {
    return this.currentUserSubject.value;
  }

  // Store tokens in localStorage
  private storeTokens(response: LoginResponse): void {
    localStorage.setItem('access_token', response.token);
    localStorage.setItem('refresh_token', response.refresh_token);
    
    const expiresAt = new Date().getTime() + (response.expires_in * 1000);
    localStorage.setItem('token_expires_at', expiresAt.toString());
  }
} 