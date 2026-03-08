using Hangfire;
using Hangfire.Dashboard;
using Hangfire.SqlServer;
using Microsoft.AspNetCore.Authentication.JwtBearer;
using Microsoft.EntityFrameworkCore;
using Microsoft.IdentityModel.Logging;
using Microsoft.IdentityModel.Tokens;
using Microsoft.OpenApi.Models;
using PortalStudioUnicaAzione.Entities;
using PortalStudioUnicaAzione.Models;
using PortalStudioUnicaAzione.Repositories;
using PortalStudioUnicaAzione.Repositories.Interfaces;
using PortalStudioUnicaAzione.Services;
using PortalStudioUnicaAzione.Services.Interfaces;
using System.Text;

var builder = WebApplication.CreateBuilder(args);
// Add JWT Authentication
var jwtSettings = builder.Configuration.GetSection("JWTSettings");
builder.Services.AddAuthentication(options =>
{
    options.DefaultAuthenticateScheme = JwtBearerDefaults.AuthenticationScheme;
    options.DefaultChallengeScheme = JwtBearerDefaults.AuthenticationScheme;
})
.AddJwtBearer(options =>
{
    options.RequireHttpsMetadata = false; // Set true in production
    options.SaveToken = true;
    options.TokenValidationParameters = new TokenValidationParameters
    {
        ValidateIssuerSigningKey = true,
        IssuerSigningKey = new SymmetricSecurityKey(Encoding.UTF8.GetBytes(jwtSettings["SecretKey"] ?? "")),
        ValidateIssuer = true,
        ValidIssuer = jwtSettings["Issuer"],
        ValidateAudience = true,
        ValidAudience = jwtSettings["Audience"],
        ValidateLifetime = false,
        ClockSkew = TimeSpan.Zero
    };
    IdentityModelEventSource.ShowPII = true;
});
builder.Services.AddAuthorization();
builder.Services.AddControllers();
builder.Services.AddHttpContextAccessor();
builder.Services.AddScoped<ICurrentUser, CurrentUser>();
// Configure Entity Framework
builder.Services.AddDbContext<PortalStudioUnicaAzioneContext>(options =>
    options.UseSqlServer(builder.Configuration.GetConnectionString("DefaultConnection")));
// Configure Swagger with JWT support
builder.Services.AddSwaggerGen(c =>
{
    c.SwaggerDoc("v1", new OpenApiInfo
    {
        Title = "Portale Studio Unica Azione API",
        Version = "v1",
        Description = "A portal to handle all our needs!"
    });
    // Add JWT authentication to Swagger
    c.AddSecurityDefinition("Bearer", new OpenApiSecurityScheme
    {
        Description = "JWT Authorization header using the Bearer scheme. Enter 'Bearer' [space] and then your token in the text input below.",
        Name = "Authorization",
        In = ParameterLocation.Header,
        Type = SecuritySchemeType.ApiKey,
        Scheme = "Bearer"
    });
    c.AddSecurityRequirement(new OpenApiSecurityRequirement()
    {
        {
            new OpenApiSecurityScheme
            {
                Reference = new OpenApiReference
                {
                    Type = ReferenceType.SecurityScheme,
                    Id = "Bearer"
                },
                Scheme = "oauth2",
                Name = "Bearer",
                In = ParameterLocation.Header,
            },
            new List<string>()
        }
    });
});
// Add distributed memory cache required for session
builder.Services.AddDistributedMemoryCache();
// Add session support for user tracking
builder.Services.AddSession(options =>
{
    options.IdleTimeout = TimeSpan.FromMinutes(30);
    options.Cookie.HttpOnly = true;
    options.Cookie.IsEssential = true;
    options.Cookie.SecurePolicy = CookieSecurePolicy.SameAsRequest;
});
// Add Hangfire services
builder.Services.AddHangfire(configuration => configuration
    .SetDataCompatibilityLevel(CompatibilityLevel.Version_180)
    .UseSimpleAssemblyNameTypeSerializer()
    .UseRecommendedSerializerSettings()
    .UseSqlServerStorage(builder.Configuration.GetConnectionString("DefaultConnection"), new SqlServerStorageOptions
    {
        CommandBatchMaxTimeout = TimeSpan.FromMinutes(5),
        SlidingInvisibilityTimeout = TimeSpan.FromMinutes(5),
        QueuePollInterval = TimeSpan.Zero,
        UseRecommendedIsolationLevel = true,
        DisableGlobalLocks = true,
        SchemaName = "HangFire" // Optional: separate schema for Hangfire tables
    }));
// Add the processing server as IHostedService
builder.Services.AddHangfireServer(options =>
{
    options.WorkerCount = Environment.ProcessorCount * 2; // Adjust based on your needs
    options.Queues = new[] { "default", "analytics" }; // Separate queue for analytics
});
// Mail Service
builder.Services.Configure<SmtpSettings>(builder.Configuration.GetSection("SmtpSettings"));
builder.Services.AddScoped<IEmailService, EmailService>();
// Repository
builder.Services.AddScoped<IAnalyticsRepository, AnalyticsRepository>();
builder.Services.AddScoped<IClickCollectionRepository, ClickCollectionRepository>();
builder.Services.AddScoped<IClickProcessingRepository, ClickProcessingRepository>();
builder.Services.AddScoped<ICompanyRepository, CompanyRepository>();
builder.Services.AddScoped<IEventsRepository, EventsRepository>();
builder.Services.AddScoped<IGeoLocationRepository, GeoLocationRepository>();
builder.Services.AddScoped<IMailTemplatesRepository, MailTemplatesRepository>();
builder.Services.AddScoped<IMediaIndexRepository, MediaIndexRepository>();
builder.Services.AddScoped<IReferrerParsingRepository, ReferrerParsingRepository>();
builder.Services.AddScoped<IShortUrlRepository, ShortUrlRepository>();
builder.Services.AddScoped<IUploadSessionRepository, UploadSessionRepository>();
builder.Services.AddScoped<IUserAgentParsingRepository, UserAgentParsingRepository>();
builder.Services.AddScoped<IUserRepository, UserRepository>();
// Services
builder.Services.AddScoped<IDbLogger, DbLogger>();
builder.Services.AddScoped<IAuthService, AuthService>();
builder.Services.AddScoped<IAnalyticsService, AnalyticsService>();
builder.Services.AddScoped<IBlobProxyService, BlobProxyService>();
builder.Services.AddSingleton<IBlobStorageService, BlobStorageService>(); // singleton ok for BlobServiceClient
builder.Services.AddScoped<IClickCollectionService, ClickCollectionService>();
builder.Services.AddScoped<IClickProcessingService, ClickProcessingService>();
builder.Services.AddScoped<ICompanyService, CompanyService>();
builder.Services.AddScoped<IEventQueryService, EventQueryService>();
builder.Services.AddScoped<IEventAccessService, EventAccessService>();
builder.Services.AddScoped<IEventsService, EventsService>();
builder.Services.AddScoped<IGeoLocationService, GeoLocationService>();
builder.Services.AddSingleton<IJwtSigningService, RsaJwtSigningService>();
builder.Services.AddScoped<IMailTemplatesService, MailTemplatesService>();
builder.Services.AddScoped<IMediaListingService, MediaListingService>();
builder.Services.AddScoped<IReferrerParsingService, ReferrerParsingService>();
builder.Services.AddScoped<IShortUrlService, ShortUrlService>();
builder.Services.AddScoped<ITokenService, TokenService>();
builder.Services.AddScoped<IUserAgentParsingService, UserAgentParsingService>();
builder.Services.AddScoped<IUploadSessionService, UploadSessionService>();
builder.Services.AddScoped<IUserService, UserService>();

// Add CORS
builder.Services.AddCors(options =>
{
    options.AddPolicy("AllowReactApp", policy =>
    {
        policy.WithOrigins("http://localhost:5173",
                                    "https://suaportal.azurewebsites.net",
                                    "https://sua.digital",
                                    "https://portal.studiounicaazione.it/")
                .AllowAnyHeader()
                .AllowAnyMethod()
                .AllowCredentials();
    });
});
var app = builder.Build();
app.UseRouting();
// Configure the HTTP request pipeline.
//if (app.Environment.IsDevelopment())
//{
app.UseDeveloperExceptionPage();
app.UseSwagger();
app.UseSwaggerUI();
//}
app.UseHangfireDashboard("/hangfire", new DashboardOptions
{
    Authorization = new[] { new RoleBasedAuthorizationFilter("Admin") }
});
app.UseHttpsRedirection();
//app.UseMiddleware<IpRateLimitingMiddleware>();
app.UseSession();
app.UseCors("AllowReactApp");
app.UseAuthentication();
app.UseAuthorization();
app.MapControllers();
// Setup recurring jobs
RecurringJob.AddOrUpdate<IClickProcessingService>(
    "process-pending-clicks",
    x => x.ProcessPendingClicksAsync(),
    "*/1 * * * *" // Every 1 minute
);
app.Run();

public class RoleBasedAuthorizationFilter : IDashboardAuthorizationFilter
{
    private readonly string _requiredRole;

    public RoleBasedAuthorizationFilter(string requiredRole)
    {
        _requiredRole = requiredRole;
    }

    public bool Authorize(DashboardContext context)
    {
        var httpContext = context.GetHttpContext();

        return httpContext.User.Identity?.IsAuthenticated == true &&
               httpContext.User.IsInRole(_requiredRole);
    }
}