import type { FieldValues, UseFormReturn } from 'react-hook-form';

export const useFormRootError = <FormValues extends FieldValues>(
  form: UseFormReturn<FormValues>
) => {
  const setRootError = (message: string) => {
    form.setError('root', {
      message,
      type: 'server'
    });
  };

  return { setRootError };
};
