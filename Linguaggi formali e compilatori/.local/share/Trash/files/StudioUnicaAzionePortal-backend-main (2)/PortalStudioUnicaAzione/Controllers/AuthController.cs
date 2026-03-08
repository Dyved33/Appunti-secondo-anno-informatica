using Microsoft.AspNetCore.Authorization;
using Microsoft.AspNetCore.Mvc;
using PortalStudioUnicaAzione.Dtos;
using PortalStudioUnicaAzione.Services.Interfaces;

namespace PortalStudioUnicaAzione.Controllers
{
    [ApiController]
    [Route("api/auth")]
    public class AuthController(IAuthService authService) : ControllerBase
    {
        [HttpPost("login")]
        public async Task<IActionResult> Login([FromBody] LoginRequest request)
        {
            try
            {
                var response = await authService.LoginAsync(request);
                return Ok(response);
            }
            catch (Exception ex)
            {
                return Problem("Login error - " + ex.Message);
            }
        }

        [HttpGet("whoami")]
        [Authorize]
        public async Task<IActionResult> WhoAmI()
        {
            try
            {
                var response = await authService.WhoAmI();
                return Ok(response);
            }
            catch (Exception ex)
            {
                return Problem("Failed to get user information - " + ex.Message);
            }
        }

        [HttpPost("forgot-password")]
        public async Task<IActionResult> ForgotPassword([FromBody] ForgotPasswordRequest request)
        {
            try
            {
                await authService.ResetPasswordRequestAsync(request.Email);
                return Ok();
            }
            catch (InvalidOperationException ex)
            {
                return Problem("Failed to generate password reset token - " + ex.Message);
            }
        }

        [HttpPost("reset-password")]
        public async Task<IActionResult> ResetPassword([FromBody] ResetPasswordRequest request)
        {
            try
            {
                await authService.ResetPasswordAsync(request.Token, request.NewPassword);
                return Ok();
            }
            catch (Exception ex)
            {
                return Problem("Failed to reset password - " + ex.Message);
            }
        }

    }
}