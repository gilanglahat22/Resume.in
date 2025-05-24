import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { AuthService } from '../../services/auth.service';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-callback',
  templateUrl: './callback.component.html',
  imports: [CommonModule],
  styleUrl: './callback.component.css'
})
export class CallbackComponent implements OnInit {
  isLoading = true;
  error = '';
  isRegistration = false;

  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private authService: AuthService
  ) {}

  ngOnInit() {
    // Check if this is a registration callback based on the URL
    const currentUrl = this.router.url;
    this.isRegistration = currentUrl.includes('/register/callback');

    // Get code and state from query parameters
    this.route.queryParams.subscribe(params => {
      const code = params['code'];
      const state = params['state'];
      const error = params['error'];

      if (error) {
        this.error = 'Authentication failed: ' + error;
        this.isLoading = false;
        return;
      }

      if (!code || !state) {
        this.error = 'Invalid callback parameters';
        this.isLoading = false;
        return;
      }

      // Handle the callback based on whether it's login or registration
      const callbackMethod = this.isRegistration 
        ? this.authService.handleRegisterCallback(code, state)
        : this.authService.handleCallback(code, state);

      callbackMethod.subscribe({
        next: (response) => {
          // Redirect to home or dashboard
          this.router.navigate(['/']);
        },
        error: (err) => {
          const action = this.isRegistration ? 'registration' : 'authentication';
          this.error = `Failed to complete ${action}`;
          this.isLoading = false;
          console.error('Callback error:', err);
        }
      });
    });
  }
} 