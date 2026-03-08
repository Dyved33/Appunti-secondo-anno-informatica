using PortalStudioUnicaAzione.Entities;

namespace PortalStudioUnicaAzione.Dtos
{
    public class UserDto
    {
        public Guid Id { get; set; } = Guid.Empty;
        public String Username { get; set; } = String.Empty;
        public String Password { get; set; } = String.Empty;
        public String Email { get; set; } = String.Empty;
        public String Firstname { get; set; } = String.Empty;
        public String Lastname { get; set; } = String.Empty;
        public String Picture { get; set; } = String.Empty;
        public Guid RoleId { get; set; } = Guid.Empty;
        public String RoleName { get; set; } = String.Empty;
        public Guid? CompanyId { get; set; }
        public String CompanyName { get; set; } = String.Empty;
        public Boolean Enabled { get; set; } = true;

        public static UserDto? ConvertToDto(Users? entity)
        {
            if (entity == null)
                return null;
            return new UserDto
            {
                Id = entity.Id,
                CompanyId = entity.CompanyId,
                Username = entity.Username,
                Email = entity.Email ?? "",
                Firstname = entity.Firstname ?? "",
                Lastname = entity.Lastname ?? "",
                Picture = entity.Picture ?? "",
                RoleId = entity.RoleId,
                RoleName = entity.Role?.Name ?? String.Empty,
                Enabled = entity.Enabled,
            };
        }

        public static UserDto? ConvertToTableDto(Users? obj)
        {
            if (obj == null) return null;
            return new UserDto()
            {
                Id = obj.Id,
                Lastname = obj.Lastname ?? "",
                Firstname = obj.Firstname ?? "",
                Username = obj.Username,
                RoleName = obj.Role.Name,
                CompanyName = obj.Company?.CompanyName ?? "",
                Enabled = obj.Enabled,
            };
        }

        public static Users? ConvertToEntity(UserDto? dto, Users entity)
        {
            if (dto == null)
                return null;
            entity.CompanyId = dto.CompanyId;
            entity.Username = dto.Username;
            if (!String.IsNullOrEmpty(dto.Password))
                entity.Password = BCrypt.Net.BCrypt.HashPassword(dto.Password);
            entity.Email = dto.Email;
            entity.Firstname = dto.Firstname;
            entity.Lastname = dto.Lastname;
            entity.RoleId = dto.RoleId;
            entity.Enabled = dto.Enabled;
            entity.Picture = dto.Picture;
            return entity;
        }
    }
}
