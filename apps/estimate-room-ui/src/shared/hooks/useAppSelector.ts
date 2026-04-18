import { useSelector } from 'react-redux';

import type { RootState } from '@/shared/types';

export const useAppSelector = useSelector.withTypes<RootState>();
