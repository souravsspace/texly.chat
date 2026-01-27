import type { ErrorComponentProps } from "@tanstack/react-router";
import {
  ErrorComponent,
  Link,
  rootRouteId,
  useMatch,
  useRouter,
} from "@tanstack/react-router";
import { AlertTriangle } from "lucide-react";

import { Button, buttonVariants } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";

export function DefaultCatchBoundary({ error }: ErrorComponentProps) {
  const router = useRouter();
  const isRoot = useMatch({
    strict: false,
    select: (state) => state.id === rootRouteId,
  });

  console.error(error);

  return (
    <div className="flex min-h-[50vh] flex-col items-center justify-center p-4">
      <Card className="w-full max-w-2xl">
        <CardHeader className="text-center">
          <div className="mb-4 flex justify-center">
            <div className="rounded-full bg-destructive/10 p-3">
              <AlertTriangle className="h-10 w-10 text-destructive" />
            </div>
          </div>
          <CardTitle className="text-xl">Something went wrong</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="max-h-[50vh] overflow-auto rounded-md bg-muted p-4 text-left">
            <ErrorComponent error={error} />
          </div>
        </CardContent>
        <CardFooter className="flex justify-center gap-2">
          <Button
            onClick={() => {
              router.invalidate();
            }}
          >
            Try Again
          </Button>
          {isRoot ? (
            <Link className={buttonVariants({ variant: "outline" })} to="/">
              Home
            </Link>
          ) : (
            <Link
              className={buttonVariants({ variant: "outline" })}
              onClick={(e) => {
                e.preventDefault();
                window.history.back();
              }}
              to="/"
            >
              Go Back
            </Link>
          )}
        </CardFooter>
      </Card>
    </div>
  );
}
