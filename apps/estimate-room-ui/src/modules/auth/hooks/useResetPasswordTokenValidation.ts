import { useCallback, useEffect, useState } from 'react';

import { ResetPasswordPageStates } from '../constants';
import {
  ResetPasswordErrorMessages,
  ResetPasswordValidationReasons
} from '../constants/resetPasswordValidation';
import { useLazyValidateResetPasswordTokenQuery } from '../store';
import type {
  ResetPasswordPageState,
  ResetPasswordValidationReason
} from '../types';
import { resolveApiErrorMessage } from '../utils/errorMessages';

const getInitialResetPasswordPageState = (token: string): ResetPasswordPageState =>
  token ? ResetPasswordPageStates.VALIDATING : ResetPasswordPageStates.INVALID;

export const useResetPasswordTokenValidation = (token: string) => {
  const [validateResetPasswordToken] = useLazyValidateResetPasswordTokenQuery();
  const [pageState, setPageState] = useState<ResetPasswordPageState>(
    getInitialResetPasswordPageState(token)
  );
  const [validationReason, setValidationReason] =
    useState<ResetPasswordValidationReason>(ResetPasswordValidationReasons.INVALID);
  const [pageError, setPageError] = useState<string | null>(null);

  const markResetLinkInvalid = useCallback(
    (
      reason: ResetPasswordValidationReason = ResetPasswordValidationReasons.INVALID,
      errorMessage: string | null = null
    ) => {
      setPageError(errorMessage);
      setValidationReason(reason);
      setPageState(ResetPasswordPageStates.INVALID);
    },
    []
  );

  useEffect(() => {
    let isMounted = true;

    const validateToken = async () => {
      if (!token) {
        markResetLinkInvalid();
        return;
      }

      setPageError(null);
      setPageState(ResetPasswordPageStates.VALIDATING);

      const response = await validateResetPasswordToken(token, true);

      if (!isMounted) {
        return;
      }

      if (response.data?.valid) {
        setPageState(ResetPasswordPageStates.READY);
        return;
      }

      if (response.data) {
        markResetLinkInvalid(response.data.reason ?? ResetPasswordValidationReasons.INVALID);
        return;
      }

      markResetLinkInvalid(
        ResetPasswordValidationReasons.INVALID,
        resolveApiErrorMessage(
          response.error,
          ResetPasswordErrorMessages.TOKEN_VALIDATION_FAILED
        )
      );
    };

    validateToken();

    return () => {
      isMounted = false;
    };
  }, [markResetLinkInvalid, token, validateResetPasswordToken]);

  return {
    markResetLinkInvalid,
    pageError,
    pageState,
    validationReason
  };
};
