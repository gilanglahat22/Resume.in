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

  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private authService: AuthService
  ) {}

  ngOnInit() {
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

      // Handle the callback (works for both login and registration flows)
      this.authService.handleCallback(code, state).subscribe({
        next: (response) => {
          // Redirect to home or dashboard
          this.router.navigate(['/']);
        },
        error: (err) => {
          // Check if it's a registration flow by state parameter
          const isRegistration = state.startsWith('register_');
          const action = isRegistration ? 'registration' : 'authentication';
          this.error = `Failed to complete ${action}`;
          this.isLoading = false;
          console.error('Callback error:', err);
        }
      });
    });
  }
} 