import { useEffect } from 'react';
import type { FieldValues, Path, UseFormReturn } from 'react-hook-form';

interface UseConfirmPasswordRevalidationOptions<FormValues extends FieldValues> {
  readonly confirmPasswordField: Path<FormValues>;
  readonly form: UseFormReturn<FormValues>;
  readonly password: string;
}

export const useConfirmPasswordRevalidation = <FormValues extends FieldValues>({
  confirmPasswordField,
  form,
  password
}: UseConfirmPasswordRevalidationOptions<FormValues>) => {
  const isConfirmPasswordTouched =
    form.getFieldState(confirmPasswordField).isTouched;

  useEffect(() => {
    if (!isConfirmPasswordTouched) {
      return;
    }

    form.trigger(confirmPasswordField);
  }, [confirmPasswordField, form, isConfirmPasswordTouched, password]);
};
