import { useDispatch } from 'react-redux';

import type { AppDispatch } from '@/shared/types';

export const useAppDispatch = useDispatch.withTypes<AppDispatch>();
