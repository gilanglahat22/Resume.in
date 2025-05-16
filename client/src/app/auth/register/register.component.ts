import { Component } from '@angular/core';
import { FormsModule } from '@angular/forms';

@Component({
  selector: 'app-register',
  imports: [FormsModule],
  templateUrl: './register.component.html',
  styleUrl: './register.component.css'
})
export class RegisterComponent {
  registerForm = { email: '', password: ''}

  onSubmit() {
    if (this.registerForm.email.trim() === '' || 
      this.registerForm.password.trim() === '') {
      return;
    }
  }
}
