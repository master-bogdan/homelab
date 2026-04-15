import {
  AppAlert,
  AppBox,
  AppButton,
  AppStack,
  AppTextField,
  AppTypography,
  SectionCard
} from '@/shared/ui';

import { useNewRoomForm } from './hooks/useNewRoomForm';
import styles from './NewRoomPage.module.scss';

export const NewRoomPage = () => {
  const {
    form: {
      formState: { errors, isSubmitting, isValid },
      register
    },
    onSubmit,
    submitMessage
  } = useNewRoomForm();

  return (
    <AppStack spacing={3}>
      <SectionCard
        description="Example React Hook Form page using typed values, MUI fields, and service-based submission."
        title="Create New Room"
      >
        <AppTypography color="text.secondary" variant="body2">
          Form state remains local to the page while backend interaction stays in the
          module service layer.
        </AppTypography>
      </SectionCard>

      <SectionCard
        description="Wire this page to the Go API once the create-room contract is available."
        title="Room Details"
      >
        <AppBox component="form" noValidate onSubmit={onSubmit}>
          <AppStack spacing={3}>
            <div className={styles.formGrid}>
              <AppTextField
                error={Boolean(errors.name)}
                helperText={errors.name?.message}
                label="Room name"
                {...register('name', {
                  minLength: {
                    message: 'Room name should be at least 2 characters.',
                    value: 2
                  },
                  required: 'Room name is required.'
                })}
              />
              <AppTextField
                error={Boolean(errors.teamId)}
                helperText={errors.teamId?.message ?? 'Optional team identifier from the backend.'}
                label="Team ID"
                {...register('teamId')}
              />
              <AppTextField
                error={Boolean(errors.length)}
                helperText={errors.length?.message}
                inputProps={{ min: 1, step: 0.1 }}
                label="Length (m)"
                type="number"
                {...register('length', {
                  min: {
                    message: 'Length must be at least 1 meter.',
                    value: 1
                  },
                  valueAsNumber: true
                })}
              />
              <AppTextField
                error={Boolean(errors.width)}
                helperText={errors.width?.message}
                inputProps={{ min: 1, step: 0.1 }}
                label="Width (m)"
                type="number"
                {...register('width', {
                  min: {
                    message: 'Width must be at least 1 meter.',
                    value: 1
                  },
                  valueAsNumber: true
                })}
              />
              <AppTextField
                error={Boolean(errors.height)}
                helperText={errors.height?.message}
                inputProps={{ min: 2, step: 0.1 }}
                label="Height (m)"
                type="number"
                {...register('height', {
                  min: {
                    message: 'Height must be at least 2 meters.',
                    value: 2
                  },
                  valueAsNumber: true
                })}
              />
            </div>

            {submitMessage ? <AppAlert severity="success">{submitMessage}</AppAlert> : null}

            <AppBox>
              <AppButton
                disabled={!isValid}
                loading={isSubmitting}
                loadingText="Creating room..."
                type="submit"
                variant="contained"
              >
                Create room scaffold
              </AppButton>
            </AppBox>
          </AppStack>
        </AppBox>
      </SectionCard>
    </AppStack>
  );
};
