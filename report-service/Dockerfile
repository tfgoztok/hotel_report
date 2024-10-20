# Use the official .NET SDK image to build the project
FROM mcr.microsoft.com/dotnet/sdk:7.0 AS build
WORKDIR /src

# Copy csproj and restore as distinct layers
COPY ["src/ReportService/ReportService.csproj", "src/ReportService/"]
RUN dotnet restore "src/ReportService/ReportService.csproj"

# Copy everything else and build
COPY src/ReportService/. src/ReportService/
WORKDIR "/src/src/ReportService"
RUN dotnet build "ReportService.csproj" -c Release -o /app/build

# Publish the application
FROM build AS publish
RUN dotnet publish "ReportService.csproj" -c Release -o /app/publish

# Build runtime image
FROM mcr.microsoft.com/dotnet/aspnet:7.0 AS final
WORKDIR /app
COPY --from=publish /app/publish .
ENTRYPOINT ["dotnet", "ReportService.dll"]