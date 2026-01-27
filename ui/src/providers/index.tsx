import { Toaster } from "@/components/ui/sonner";
import { AlertDialogProvider } from "@/providers/alert-dialog";
import { AuthProvider } from "@/providers/auth";
import { ThemeProvider } from "@/providers/theme";

const Providers = ({ children }: { children: React.ReactNode }) => {
  return (
    <AuthProvider>
      <ThemeProvider>
        <AlertDialogProvider>{children}</AlertDialogProvider>
        <Toaster />
      </ThemeProvider>
    </AuthProvider>
  );
};

export default Providers;
