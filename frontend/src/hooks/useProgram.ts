import { useQuery } from '@tanstack/react-query';
import { queryKeys } from '../lib/query';
import { programs } from '../api';

export function useProgram(id: string | undefined) {
  return useQuery({
    queryKey: queryKeys.programs.detail(id!),
    queryFn: () => programs.getProgram(id!),
    enabled: !!id,
  });
}

export function useProgramBySlug(slug: string | undefined) {
  return useQuery({
    queryKey: queryKeys.programs.bySlug(slug!),
    queryFn: () => programs.getProgramBySlug(slug!),
    enabled: !!slug,
  });
}
