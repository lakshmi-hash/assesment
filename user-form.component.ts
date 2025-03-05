import { Component, Inject } from '@angular/core';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';
import { UserService } from '../../services/user.service';

@Component({
  selector: 'app-user-form',
  templateUrl: './user-form.component.html',
  styleUrls: ['./user-form.component.css']
})
export class UserFormComponent {
  user: any = { name: '', email: '' };
  error: string = '';

  constructor(
    public dialogRef: MatDialogRef<UserFormComponent>,
    @Inject(MAT_DIALOG_DATA) public data: any,
    private userService: UserService
  ) {
    if (data.id) {
      this.user = { ...data };
    }
  }

  save() {
    if (!this.user.name || !this.user.email || !this.isValidEmail(this.user.email)) {
      this.error = 'Name and valid email are required';
      return;
    }

    if (this.user.id) {
      this.userService.updateUser(this.user.id, this.user).subscribe({
        next: () => this.dialogRef.close(true),
        error: (err) => this.error = err.error.error || 'Failed to update user'
      });
    } else {
      this.userService.createUser(this.user).subscribe({
        next: () => this.dialogRef.close(true),
        error: (err) => this.error = err.error.error || 'Failed to create user'
      });
    }
  }

  isValidEmail(email: string): boolean {
    const re = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    return re.test(email);
  }
}