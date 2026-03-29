import { SignupForm } from "./components";


export default function SetupPage() {
  return (
    <div className="relative grid min-h-svh w-screen max-h-screen overflow-clip lg:grid-cols-2">
      <div className="flex flex-col z-20 gap-4 p-6 md:p-10 w-screen justify-center min-h-svh backdrop-blur-sm lg:backdrop-blur-none">
        <div className="flex justify-center gap-2 md:justify-start">
          <div className="flex items-center gap-2">
            <div className="flex items-center justify-center size-10 rounded-lg bg-primary/10 border border-primary/20">
              <img src="/mist.png" alt="Mist Logo" className="size-6" />
            </div>
            <span className="font-bold text-2xl tracking-tight">Mist</span>
          </div>
        </div>
        <div className="flex flex-1 mx-auto items-center w-full justify-center">
          <div className="w-full max-w-md">
            <div className="rounded-xl border bg-card/50 backdrop-blur-md p-8 shadow-lg">
              <SignupForm />
            </div>
          </div>
        </div>
      </div>
      <div className="hidden lg:block relative">
        <div className="absolute inset-0 bg-gradient-to-br from-primary/20 via-background to-background z-10" />
        <img
          src="/cloud-computing.png"
          alt="Cloud Computing"
          className="absolute inset-0 h-full w-full object-cover dark:brightness-[0.4] dark:grayscale"
        />
      </div>
      <img
        src="/cloud-computing.png"
        alt="Background"
        className="absolute inset-0 h-full w-full object-cover dark:brightness-[0.2] dark:grayscale lg:hidden -z-10"
      />
    </div>
  )
}
