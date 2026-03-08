using PortalStudioUnicaAzione.Entities;
using static Amazon.S3.Util.S3EventNotification;

namespace PortalStudioUnicaAzione.Dtos
{
    public class CompanyDto
    {
        public Guid Id { get; set; } = Guid.Empty;
        public string CompanyName { get; set; } = string.Empty;
        public string MainContactPerson { get; set; } = string.Empty;
        public string? Email { get; set; }
        public string? Phone { get; set; }
        public bool Enabled { get; set; } = true;

        public static CompanyDto? ConvertToDto(Companies? entity)
        {
            if (entity == null)
                return null;
            return new CompanyDto
            {
                Id = entity.Id,
                CompanyName = entity.CompanyName ?? "",
                MainContactPerson = entity.MainContactPerson ?? "",
                Email = entity.Email,
                Phone = entity.Phone,
                Enabled = entity.Enabled,
            };
        }

        public static CompanyDto? ConvertToTableDto(Companies? entity)
        {
            if (entity == null) return null;
            return new CompanyDto()
            {
                Id = entity.Id,
                CompanyName = entity.CompanyName ?? "",
                MainContactPerson = entity.MainContactPerson ?? "",
                Email = entity.Email,
                Phone = entity.Phone,
                Enabled = entity.Enabled,
            };
        }

        public static Companies? ConvertToEntity(CompanyDto? dto, Companies entity)
        {
            if (dto == null)
                return null;
            entity.CompanyName = dto.CompanyName;
            entity.MainContactPerson = dto.MainContactPerson;
            entity.Email = dto.Email;
            entity.Phone = dto.Phone;
            entity.Enabled = dto.Enabled;
            return entity;
        }
    }
}