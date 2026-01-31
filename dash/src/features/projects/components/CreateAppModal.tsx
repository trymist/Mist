import { useState } from "react";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { AppTypeSelection } from "./AppTypeSelection";
import { WebAppForm } from "./WebAppForm";
import { ServiceForm } from "./ServiceForm";
import { DatabaseForm } from "./DatabaseForm";
import { ComposeAppForm } from "./ComposeAppForm";
import type { AppType, CreateAppRequest } from "@/types/app";

interface CreateAppModalProps {
  isOpen: boolean;
  onClose: () => void;
  projectId: number;
  onSubmit: (data: CreateAppRequest) => void | Promise<void>;
}

type Step = "select-type" | "configure";

export function CreateAppModal({
  isOpen,
  onClose,
  projectId,
  onSubmit,
}: CreateAppModalProps) {
  const [step, setStep] = useState<Step>("select-type");
  const [selectedType, setSelectedType] = useState<AppType | null>(null);

  const handleTypeSelect = (type: AppType) => {
    setSelectedType(type);
    setStep("configure");
  };

  const handleBack = () => {
    setStep("select-type");
    setSelectedType(null);
  };

  const handleSubmit = async (data: CreateAppRequest) => {
    await onSubmit(data);
    handleClose();
  };

  const handleClose = () => {
    setStep("select-type");
    setSelectedType(null);
    onClose();
  };

  return (
    <Dialog open={isOpen} onOpenChange={handleClose}>
      <DialogContent className="sm:max-w-2xl max-h-[85vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle>
            {step === "select-type" ? "Create New Application" : "Configure Application"}
          </DialogTitle>
        </DialogHeader>

        {step === "select-type" && (
          <AppTypeSelection onSelect={handleTypeSelect} />
        )}

        {step === "configure" && selectedType === "web" && (
          <WebAppForm
            projectId={projectId}
            onSubmit={handleSubmit}
            onBack={handleBack}
          />
        )}

        {step === "configure" && selectedType === "service" && (
          <ServiceForm
            projectId={projectId}
            onSubmit={handleSubmit}
            onBack={handleBack}
          />
        )}

        {step === "configure" && selectedType === "database" && (
          <DatabaseForm
            projectId={projectId}
            onSubmit={handleSubmit}
            onBack={handleBack}
          />
        )}

        {step === "configure" && selectedType === "compose" && (
          <ComposeAppForm
            projectId={projectId}
            onSubmit={handleSubmit}
            onBack={handleBack}
          />
        )}
      </DialogContent>
    </Dialog>
  );
}
