import { Link } from "@tanstack/react-router";
import { FileQuestion } from "lucide-react";

import { Button, buttonVariants } from "@/components/ui/button";
import {
  Card,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";

export function NotFound({ children }: { children?: React.ReactNode }) {
  return (
    <div className="flex min-h-[50vh] flex-col items-center justify-center p-4">
      <Card className="w-full max-w-md text-center">
        <CardHeader>
          <div className="mb-4 flex justify-center">
            <div className="rounded-full bg-muted p-3">
              <FileQuestion className="h-10 w-10 text-muted-foreground" />
            </div>
          </div>
          <CardTitle className="text-2xl">Page Not Found</CardTitle>
          <CardDescription>
            {children || "The page you are looking for does not exist."}
          </CardDescription>
        </CardHeader>
        <CardFooter className="flex justify-center gap-2">
          <Button onClick={() => window.history.back()} variant="outline">
            Go Back
          </Button>
          <Link className={buttonVariants({ variant: "default" })} to="/">
            Start Over
          </Link>
        </CardFooter>
      </Card>
    </div>
  );
}
