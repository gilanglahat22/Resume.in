import { Component } from '@angular/core';
import { FormsModule } from '@angular/forms';

@Component({
  selector: 'app-login',
  imports: [FormsModule],
  templateUrl: './login.component.html',
  styleUrl: './login.component.css'
})
export class LoginComponent {
  loginForm = { email: '', password: '' };

  onSubmit() {
    if (this.loginForm.email.trim() === '' || 
      this.loginForm.password.trim() === '') {
      return;
    }
  }
}
