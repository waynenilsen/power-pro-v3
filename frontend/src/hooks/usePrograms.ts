import { useQuery } from '@tanstack/react-query';
import { queryKeys } from '../lib/query';
import { programs } from '../api';
import type { ListProgramsParams } from '../api/endpoints/programs';

export function usePrograms(params?: ListProgramsParams) {
  return useQuery({
    queryKey: queryKeys.programs.list(params),
    queryFn: () => programs.listPrograms(params),
  });
}
