import { create } from 'zustand';
import { devtools } from 'zustand/middleware';

interface HistoryStack {
  stack: string[];
  push: (path: string) => void;
  pop: () => void;
  replace: (path: string) => void;
}

const useHistoryStackStore = create<HistoryStack>()(
  devtools((set) => ({
    stack: [],
    push: (path) => set((prev) => ({ stack: [...prev.stack, path] })),
    pop: () =>
      set((prev) => ({ stack: prev.stack.slice(0, prev.stack.length - 1) })),
    replace: (path) =>
      set((prev) => ({
        stack: [...prev.stack.slice(0, prev.stack.length - 1), path],
      })),
  })),
);

export default useHistoryStackStore;