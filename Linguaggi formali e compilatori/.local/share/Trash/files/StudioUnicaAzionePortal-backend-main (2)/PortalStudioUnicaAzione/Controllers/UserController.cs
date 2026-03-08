using Microsoft.AspNetCore.Authorization;
using Microsoft.AspNetCore.Mvc;
using PortalStudioUnicaAzione.Dtos;
using PortalStudioUnicaAzione.Models;
using PortalStudioUnicaAzione.Services.Interfaces;

namespace PortalStudioUnicaAzione.Controllers
{
    [ApiController]
    [Route("api/[controller]")]
    [Authorize(Roles = "Admin")]
    public class UserController(IUserService userService) : ControllerBase
    {
        [HttpPost("list")]
        public async Task<ActionResult<TablesResponse<UserDto?>>> GetTableList(TablesUsersRequest request)
        {
            try
            {
                if (!ModelState.IsValid)
                {
                    var errors = (from item in ModelState where item.Value.Errors.Any() select item.Value.Errors[0].ErrorMessage).ToList();
                    return BadRequest(new CodeLevelResponse(MessageType.Error, "Problems with received data!", String.Join(";", [.. errors])));
                }
                var users = await userService.GetTableList(request);
                return Ok(users);
            }
            catch (Exception ex)
            {
                return BadRequest(new { message = ex.Message });
            }
        }

        [HttpGet("{id}")]
        public async Task<ActionResult<UserDto>> Get(Guid id)
        {
            try
            {
                var user = await userService.GetById(id);
                return Ok(user);
            }
            catch (UnauthorizedAccessException ex)
            {
                return Forbid(ex.Message);
            }
            catch (KeyNotFoundException ex)
            {
                return NotFound(new { message = ex.Message });
            }
            catch (Exception ex)
            {
                return BadRequest(new { message = ex.Message });
            }
        }

        [HttpPost]
        public async Task<ActionResult> Insert([FromBody] UserDto request)
        {
            try
            {
                var entity = UserDto.ConvertToEntity(request, new Entities.Users());
                var user = await userService.Insert(entity);
                return CreatedAtAction(nameof(Insert), new { id = user?.Id });
            }
            catch (UnauthorizedAccessException ex)
            {
                return Forbid(ex.Message);
            }
            catch (Exception ex)
            {
                return BadRequest(new { message = ex.Message });
            }
        }

        [HttpPut]
        public async Task<ActionResult> Update([FromBody] UserDto request)
        {
            try
            {
                var existing = await userService.GetById(request.Id);
                if (existing == null)
                    return NotFound();
                var entity = UserDto.ConvertToEntity(request, existing);
                await userService.Update(entity);
                return Ok();
            }
            catch (UnauthorizedAccessException)
            {
                return Forbid();
            }
            catch (KeyNotFoundException)
            {
                return NotFound();
            }
            catch (Exception)
            {
                return BadRequest();
            }
        }

        [HttpDelete("{id}")]
        public async Task<IActionResult> SafeDelete(Guid id)
        {
            try
            {
                await userService.SafeDelete(id);
                return NoContent();
            }
            catch (UnauthorizedAccessException)
            {
                return Forbid();
            }
            catch (KeyNotFoundException)
            {
                return NotFound();
            }
            catch (Exception)
            {
                return BadRequest();
            }
        }

        [HttpGet("{id}/enable")]
        public async Task<IActionResult> Enable(Guid id)
        {
            await userService.Enable(id);
            return Ok();
        }

        [HttpGet("{id}/disable")]
        public async Task<IActionResult> Disable(Guid id)
        {
            await userService.Disable(id);
            return Ok();
        }
    }
}