import { Component } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { CommonModule } from '@angular/common';
import { AuthService } from '../../services/auth.service';
import { Router } from '@angular/router';

@Component({
  selector: 'app-login',
  imports: [FormsModule, CommonModule],
  templateUrl: './login.component.html',
  styleUrl: './login.component.css'
})
export class LoginComponent {
  loginForm = { email: '', password: '' };
  isLoading = false;
  errorMessage = '';

  constructor(
    private authService: AuthService,
    private router: Router
  ) {}

  onSubmit() {
    if (this.loginForm.email.trim() === '' || 
      this.loginForm.password.trim() === '') {
      this.errorMessage = 'Please fill in all fields';
      return;
    }
    
    // For now, show a message that only OAuth login is supported
    this.errorMessage = 'Please use Google login for authentication';
  }

  googleLogin() {
    this.isLoading = true;
    this.errorMessage = '';
    
    this.authService.googleLogin().subscribe({
      next: (response) => {
        // Redirect to Google OAuth
        window.location.href = response.auth_url;
      },
      error: (error) => {
        this.isLoading = false;
        this.errorMessage = 'Failed to initiate Google login. Please try again.';
        console.error('Google login error:', error);
      }
    });
  }
}
