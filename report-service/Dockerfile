# Build stage
FROM mcr.microsoft.com/dotnet/sdk:7.0 AS build
WORKDIR /src

# Copy csproj and restore as distinct layers
COPY ["src/ReportService/ReportService.csproj", "ReportService/"]
RUN dotnet restore "ReportService/ReportService.csproj"

# Copy everything else and build
COPY src/ReportService/. ReportService/
WORKDIR "/src/ReportService"
RUN dotnet build "ReportService.csproj" -c Release -o /app/build

# Test stage
FROM build AS test
WORKDIR /src
COPY ["tests/ReportService.Tests/ReportService.Tests.csproj", "ReportService.Tests/"]
RUN dotnet restore "ReportService.Tests/ReportService.Tests.csproj"
COPY tests/ReportService.Tests/. ReportService.Tests/
WORKDIR "/src/ReportService.Tests"
RUN dotnet build "ReportService.Tests.csproj" -c Release -o /app/build
RUN dotnet test "ReportService.Tests.csproj" --logger:trx

# Publish stage
FROM build AS publish
WORKDIR "/src/ReportService"
RUN dotnet publish "ReportService.csproj" -c Release -o /app/publish

# Final stage
FROM mcr.microsoft.com/dotnet/aspnet:7.0 AS final
WORKDIR /app
COPY --from=publish /app/publish .
ENTRYPOINT ["dotnet", "ReportService.dll"]